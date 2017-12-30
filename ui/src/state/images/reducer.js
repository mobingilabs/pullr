import * as Actions from './actions';

export default function (state, action) {
    switch (action.type) {
        case Actions.IMAGES_SHOW_PAGE:
            return showPage(state, action);
        case Actions.IMAGES_ADD_IMAGE:
            return addImage(state, action);
        case Actions.IMAGES_UPDATE_IMAGE:
            return updateImage(state, action);
        default:
            return state;
    }
}

export function showPage(state, action) {
    return Object.assign({}, state, {
        currentPage: action.pageNumber
    });
}

export function addImage(state, action) {
    const dataOrder = [].concat(state.dataOrder, [action.image.name]);
    const data = Object.assign({}, state.data, { [action.image.name]: action.image });
    return Object.assign({}, state, { dataOrder, data });
}

export function updateImage(state, action) {
    let dataOrder = state.dataOrder;
    let data = Object.assign({}, state.data);
    if (action.imageName != action.imageData.name) {
        dataOrder = [].concat(dataOrder);

        const oldNameIndex = dataOrder.indexOf(action.imageName);
        dataOrder[oldNameIndex] = action.imageData.name;

        delete data[action.imageName];
    }

    data[action.imageData.name] = action.imageData;

    return Object.assign({}, state, { dataOrder, data });
}
