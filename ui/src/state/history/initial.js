export default {
    currentPage: 0,
    loadInProgress: false,
    thereIsMore: true,
    lastBuilds: [
        {
            imageName: 'pullr',
            date: new Date(),
            tag: 'latest',
            status: 'inProgress'
        }
    ],
    buildsByImage: {
        pullr: {
            image: 'pullr',
            build: {
                type: 'branch',
                name: '/v.*/'
            },
            gitRef: {

            }
        }
    }
}
