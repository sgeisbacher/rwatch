// TODO deduplicate
export interface ExecutionInfo {
	command: string;
	exec_time: Date;
	exec_count: number;
	success: boolean;
	output: string;
}
