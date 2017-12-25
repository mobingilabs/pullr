export default {
    currentPage: 0,
    totalPages: 20,
    totalImages: 200,
    details: {
        'pullr-icarium': {
            name: 'pullr-icarium',
            provider: 'github',
            organisation: 'mobingilabs',
            repository: 'pullr-icarium',
            dockerfilePath: 'Dockerfile',
            tags: [
                { type: 'branch', name: 'master', tag: 'latest' },
                { type: 'tag', name: '/v.*/' }
            ],
        },
        'pullr': {
            name: 'pullr',
            provider: 'github',
            organisation: 'mobingilabs',
            repository: 'pullr-icarium',
            dockerfilePath: 'Dockerfile',
            tags: [
                { type: 'branch', name: 'master', tag: 'latest' },
                { type: 'tag', name: '/v.*/' }
            ],
        }
    },
    data: [
        {
            name: 'pullr-icarium',
            commitHash: 'dc6a977bbc80ea581ce7f08362822c9650caa7d2',
            tag: 'v0.0.1',
            lastBuild: new Date()
        },
        {
            name: 'pullr',
            commitHash: 'dc6a977bbc80ea581ce7f08362822c9650caa7d2',
            tag: 'v0.0.1',
            lastBuild: new Date()
        }
    ]
}
