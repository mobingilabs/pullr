export default class SourceService {
  constructor (http) {
    this.http = http
  }

  /**
   * owners, load source owners for given source provider
   * @param {string} provider
   * @returns {Promise<string[]>}
   */
  owners (provider) {
    return this.http.get(`/source/${provider}/orgs`).then(res => res.data)
  }

  /**
   * repositories, load source repositories for given owner
   * @param {string} provider
   * @param {string} owner
   * @returns {Promise<string[]>}
   */
  repositories (provider, owner) {
    const params = {org: owner}
    return this.http.get(`/source/${provider}/repos`, {params}).then(res => res.data)
  }
}
