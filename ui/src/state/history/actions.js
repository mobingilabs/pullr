export const HISTORY_ADD_PAGE = 'HISTORY_ADD_PAGE';

export function addPage(builds) {
    return {
        type: HISTORY_ADD_PAGE,
        builds
    };
}

export function loadHistory(page) {
    return (dispatch) => {
        
    };
}
