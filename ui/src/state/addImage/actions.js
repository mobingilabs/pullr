import { addImage } from '../images/actions';

export const ADD_IMAGE_SET_PROVIDER = 'ADD_IMAGE_SET_PROVIDER';
export const ADD_IMAGE_SET_ORGANISATION = 'ADD_IMAGE_SET_ORGANISATION';
export const ADD_IMAGE_SET_REPOSITORY = 'ADD_IMAGE_SET_REPOSITORY';
export const ADD_IMAGE_UPDATE_CONFIG = 'ADD_IMAGE_UPDATE_CONFIG';
export const ADD_IMAGE_ADD_BUILD_TAG = 'ADD_IMAGE_ADD_BUILD_TAG';
export const ADD_IMAGE_RESET = 'ADD_IMAGE_RESET';

export function setProvider(provider) {
    return {
        type: ADD_IMAGE_SET_PROVIDER,
        provider
    };
}

export function setOrganisation(organisation) {
    return {
        type: ADD_IMAGE_SET_ORGANISATION,
        organisation
    };
}

export function setRepository(repository) {
    return {
        type: ADD_IMAGE_SET_REPOSITORY,
        repository
    };
}

export function reset() {
    return {
        type: ADD_IMAGE_RESET
    };
}

export function updateConfig(config) {
    return {
        type: ADD_IMAGE_UPDATE_CONFIG,
        config
    };
}

export function createImage(image) {
    return (dispatch) => {

    };
}
