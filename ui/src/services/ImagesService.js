export default class ImagesService {
  /**
   * @param {AxiosInstance} http
   */
  constructor (http) {
    this.http = http

    // Caching images request promises to avoid concurrent
    // fetching of images list
    this._imagesPromise = null
  }

  /**
   * Images fetches images from the api server. Fetched data
   * contains both array of images and the pagination information
   * @returns {Promise<{images: Object[], pagination: Object}[]>}
   */
  list () {
    if (this._imagesPromise) {
      return this._imagesPromise
    }

    this._imagesPromise = this.http.get('/images').then(res => res.data)
    this._imagesPromise.catch(() => {}).then(() => { this._imagesPromise = null })
    return this._imagesPromise
  }

  /**
   * updateImage updates the image data on the api server
   * @param {string} key Image key
   * @param {Object} data Image data
   * @returns {Promise<Object>}
   */
  update (key, data) {
    return this.http.post(`/images/${key}`, data).then(res => res.data)
  }

  /**
   * deleteImage deletes an image on the api server
   * @param {string} key Image key
   * @returns {Promise<null>}
   */
  delete (key) {
    return this.http.delete(`/images/${key}`).then(() => null)
  }

  /**
   * createImage stores the image information on the api server
   * @param {Object} data Image data
   * @returns {Promise<null>}
   */
  create (data) {
    return this.http.post('/images', data).then(() => null)
  }
}
