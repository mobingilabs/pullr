import * as actions from './actions';

export default function (state, action) {
    switch(action.type) {
        case actions.HISTORY_START_LOADING:
            return startLoading(state, action);
        case actions.HISTORY_FINISH_LOADING:
            return finishLoading(state, action);
        default:
            return state;
    }
}

function startLoading(state, action) {
    return {
        ...state,
        loadInProgress: true
    };
}

function finishLoading(state, action) {
    return {
        ...state,
        loadInProgress: false,
        currentPage: action.pageNumber,
        lastBuilds: [].concat(state.lastBuilds, action.builds),
        thereIsMore: action.builds.length > 0
    };
}
