export const HISTORY_ADD_PAGE = 'HISTORY_ADD_PAGE';
export const HISTORY_START_LOADING = 'HISTORY_START_LOADING';
export const HISTORY_FINISH_LOADING = 'HISTORY_FINISH_LOADING';

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

export function startLoading(pageNumber) {
    return {
        type: HISTORY_START_LOADING,
        pageNumber
    };
}

export function finishLoading(pageNumber, builds, error) {
    return {
        type: HISTORY_FINISH_LOADING,
        pageNumber,
        builds,
        error
    };
}

export function loadMore() {
    return (dispatch, getState) => {
        const { currentPage } = getState().history;
        dispatch(startLoading(currentPage + 1));
        setTimeout(() => {
            dispatch(finishLoading(currentPage + 1, []));
        }, 500);
    };
}
