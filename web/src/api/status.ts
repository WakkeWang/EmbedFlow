import type { TagProps } from 'naive-ui'

// Shared status -> tag-type mapping (batch statuses and build-record
// statuses share the vocabulary). One map instead of a per-view switch
// (the review's Repeated Switches finding).
export function statusType(
	s: string,
): TagProps['type'] {
	switch (s) {
		case 'queued':
		case 'pending':
		case 'closed':
			return 'default'
		case 'running':
		case 'building':
		case 'active':
			return 'info'
		case 'completed':
		case 'succeeded':
			return 'success'
		case 'failed':
			return 'error'
		case 'canceled':
		case 'skipped':
			return 'warning'
		default:
			return 'default'
	}
}
