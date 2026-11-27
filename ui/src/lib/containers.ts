import type { Fn, Status } from '#lib/types.js';

export interface Container {
	Id: string;
	Names: string[] | null;
	Image: string;
	State: string; // created, running, paused, exited, stopped, dead
	Created: string; // ISO date
	Ports: { host_ip: string; container_port: number; host_port: number; protocol: string }[] | null;
}

// What one row of the Functions table shows.
export type Row = Pick<Fn, 'name' | 'status' | 'image' | 'address' | 'requests24h' | 'lastDeploy'>;

function toStatus(state: string): Status {
	if (state === 'running') return 'Running';
	if (state === 'dead') return 'Failed';
	return 'Stopped';
}

function ago(iso: string): string {
	const s = Math.round((Date.now() - new Date(iso).getTime()) / 1000);
	if (s < 60) return 'just now';
	if (s < 3600) return `${Math.floor(s / 60)} min ago`;
	if (s < 86400) return `${Math.floor(s / 3600)} h ago`;
	return `${Math.floor(s / 86400)} d ago`;
}

export function toRow(c: Container): Row {
	const port = c.Ports?.find((p) => p.host_port > 0);
	return {
		name: c.Names?.[0] ?? c.Id.slice(0, 12),
		status: toStatus(c.State),
		image: c.Image,
		address: c.State === 'running' && port ? `localhost:${port.host_port}` : null,
		requests24h: null,
		lastDeploy: ago(c.Created)
	};
}
