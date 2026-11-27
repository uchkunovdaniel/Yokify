import { toRow, type Container, type Row } from '#lib/containers.js';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
	let rows: Row[] = [];
	let error: string | null = null;
	try {
		const res = await fetch('/api/containers');
		if (res.ok) {
			const list: Container[] | null = await res.json();
			rows = (list ?? []).map(toRow);
		} else {
			error = (await res.text()).trim() || `The server answered ${res.status}.`;
		}
	} catch {
		error = 'The Yōkify server is not running.';
	}
	return { rows, error };
};
