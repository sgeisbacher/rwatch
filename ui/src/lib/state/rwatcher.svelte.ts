export interface LogEvent {
	level: 'debug' | 'info' | 'warn' | 'error';
	msg: string;
}
export interface RwatcherState {
	logs: LogEvent[];
}

// export const rwatcher = $state<RwatcherState>({ logs: [] });
// setInterval(() => console.log(' indexxxxxx '), 2000);
// console.log('lllllllllll');
