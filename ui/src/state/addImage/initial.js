import * as helpers from './helpers';

export default {
    step: 0,
    organisations: ['mobingilabs', 'EpicGames', 'umurgdk'],
    repositories: {
        mobingilabs: ['pullr'],
        umurgdk: ['soundlines', 'flight-crusader']
    },
    image: {
        provider: null,
        organisation: null,
        repository: null,
        name: '',
        dockerfilePath: '',
        builds: [helpers.defaultBuildObject()],
    }
}
