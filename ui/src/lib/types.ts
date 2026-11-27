export type Status = 'Running' | 'Building' | 'Failed' | 'Stopped';

export type Source =
	| { kind: 'git'; provider: 'GitHub' | 'GitLab' | 'Bitbucket'; repo: string; branch: string }
	| { kind: 'folder'; path: string };

export interface Fn {
	name: string;
	status: Status;
	image: string;
	address: string | null;
	requests24h: number | null;
	lastDeploy: string;
	source: Source;
	runtime: string;
}

export interface Deployment {
	fn: string;
	commit: string;
	message: string;
	author: string;
	when: string;
	duration: string;
	result: 'Success' | 'Failed' | 'Building';
	live?: boolean;
}

export type Level = 'info' | 'warn' | 'error';
export interface LogLine {
	t: string;
	level: Level;
	text: string;
}

export function sourceLabel(s: Source): string {
	return s.kind === 'git' ? `${s.provider} · ${s.repo} · ${s.branch}` : s.path;
}
