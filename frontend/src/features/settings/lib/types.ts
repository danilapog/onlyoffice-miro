export interface SettingsRequest {
    address: string;
    header: string;
    secret: string;
    demo: boolean;
}

export interface SettingsResponse {
    address: string;
    header: string;
    secret: string;
    demo: {
        team_id: string;
        enabled: boolean;
        started: string;
    };
}