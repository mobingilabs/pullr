import * as Actions from './actions';

export function showPage(state, action) {
    return Object.assign({}, state, {
        currentPage: action.pageNumber
    });
}

export function addImage(state, action) {
    return Object.assign({}, state, {
        data: state.images.concat([ action.image ])
    });
}

export function addDetail(state, action) {
    const details = Object.assign({}, state.details);
    details[action.imageDetail.name] = action.imageDetail;

    return Object.assign({}, state, { details });
}

export default function (state, action) {
    switch (action.type) {
        case Actions.IMAGES_SHOW_PAGE:
            return showPage(state, action);
        case Actions.IMAGES_ADD_ENTRY:
            return addImage(state, action);
        case Actions.IMAGES_ADD_DETAIL:
            return addDetail(state, action);
        default:
            return state;
    }
}
