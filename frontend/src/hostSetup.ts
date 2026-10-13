// The host screen's remembered setup (UC-03), kept in the settings.

// HostSetup is a remembered host screen; seats 0 and sets null mean
// "not set yet". The password is never kept (N-14).
export interface HostSetup {
  seats: number;
  sets: string[] | null;
  detailed: boolean;
  response_timer: number;
  no_bonus_souls: boolean;
}

// noHostSetup is the setup before the first host.
export const noHostSetup: HostSetup = { seats: 0, sets: null, detailed: false, response_timer: 0, no_bonus_souls: false };
