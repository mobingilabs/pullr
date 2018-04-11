export default class BuildService {
  /**
   * @param {AxiosInstance} http
   */
  constructor (http) {
    this.http = http
  }

  /**
   * list, lists recent builds unique by image
   * @returns {Promise<{builds: Array, images: Object}>}
   */
  list () {
    return this.http.get('/builds').then(res => res.data)
  }

  /**
   * history, lists the build history for the given image key
   * @param {string} imgKey
   * @param {{page: number, per_page?: number}} pagination
   * @returns {Promise<{build_records: Array, pagination: Object}>}
   */
  history (imgKey, pagination) {
    if (!pagination) {
      pagination = {page: 0}
    }
    return this.http.get(`/builds/${imgKey}`, {params: pagination}).then(res => res.data)
  }
}
