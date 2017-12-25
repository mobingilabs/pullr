import * as helpers from './helpers';

export default {
    step: 0,
    provider: null,
    organisation: null,
    repository: null,
    organisations: ['mobingilabs', 'EpicGames', 'umurgdk'],
    repositories: {
        mobingilabs: ['pullr'],
        umurgdk: ['soundlines', 'flight-crusader']
    },
    config: {
        dockerfilePath: '',
        tags: [helpers.defaultTagObject()],
    }
}
