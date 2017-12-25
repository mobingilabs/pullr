import * as actions from './actions';
import initialState from './initial';

export default function (state, action) {
    switch (action.type) {
        case actions.ADD_IMAGE_SET_PROVIDER:
            return setProvider(state, action);
        case actions.ADD_IMAGE_SET_ORGANISATION:
            return setOrganisation(state, action);
        case actions.ADD_IMAGE_SET_REPOSITORY:
            return setRepository(state, action);
        case actions.ADD_IMAGE_UPDATE_CONFIG:
            return updateConfig(state, action);
        case actions.ADD_IMAGE_RESET:
            return reset(state, action);
        default:
            return state;
    }
}

export function setProvider(state, action) {
    return Object.assign({}, state, {
        step: 1,
        provider: action.provider
    });
}

export function setOrganisation(state, action) {
    return Object.assign({}, state, {
        organisation: action.organisation
    });
}

export function setRepository(state, action) {
    return Object.assign({}, state, {
        step: 2,
        repository: action.repository
    });
}

export function updateConfig(state, action) {
    return Object.assign({}, state, { config: action.config });
}

export function reset(state, action) {
    return Object.assign({}, initialState, {
        repositories: state.repositories,
        organisations: state.organisations
    });
}
