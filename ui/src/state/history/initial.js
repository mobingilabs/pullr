export default {
    page: 0,
    builds: [
        {
            image: 'pullr',
            build: {
                type: 'branch',
                name: '/v.*/'
            },
            gitRef: {

            }
        }
    ]
}
