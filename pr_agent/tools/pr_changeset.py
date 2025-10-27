import re
import traceback

from pr_agent.config_loader import get_settings
from pr_agent.git_providers import GitLabProvider
from pr_agent.log import get_logger

def find_changesets(text):
    # Regular expression patterns for Changesets
    patterns = [
        r'https://gitlab\.urent\.city/urent/urentbike-doc/-/merge_requests/[0-9]+',
        r'https://gitlab\.urent\.city/urent/urentbike-doc/-/.+\.md'
    ]

    # Regular expression patterns for Changesets
    changesets = set()
    for pattern in patterns:
        matches = re.findall(pattern, text)
        for match in matches:
            if match:
                changesets.add(match)

    return list(changesets)

async def extract_changeset(git_provider):
    try:
        if isinstance(git_provider, GitLabProvider):
            user_description = git_provider.get_user_description()
            changeset = find_changesets(user_description)[0]
            if ".md" in changeset:
                return git_provider.get_changeset_content_from_file(changeset)
            else:
                return git_provider.get_changeset_content_from_mr(changeset)
    except Exception as e:
        get_logger().warning(f"Error extracting changeset error= {e}",
                            artifact={"traceback": traceback.format_exc()})
        return None

async def extract_and_cache_changeset(git_provider, vars):
    if not get_settings().get('pr_reviewer.require_changeset_review', False):
        return

    cached_changeset = get_settings().get('changeset', None)

    if not cached_changeset:
        cached_changeset = await extract_changeset(git_provider)
        if cached_changeset:
            get_logger().info("Extracted changeset from GitLab",
                        artifact={"changeset": cached_changeset})

            vars['changeset'] = cached_changeset
            get_settings().set('changeset', cached_changeset)
    else:
        get_logger().info(f"Using cached changeset: {cached_changeset}", artifact={"changeset": cached_changeset})
        vars['changeset'] = cached_changeset

