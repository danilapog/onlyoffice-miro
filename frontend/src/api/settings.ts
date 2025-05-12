export const checkSettings: () => Promise<boolean> = async () => {
    const { board: miroBoard } = window.miro;
    const settings = await miroBoard.getAppData("settings");
    return !!settings;
};