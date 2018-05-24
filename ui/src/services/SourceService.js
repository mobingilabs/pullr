export default class SourceService {
  constructor (http) {
    this.http = http
  }

  /**
   * organisations, load source organisations for given source provider
   * @param {string} provider
   * @returns {Promise<string[]>}
   */
  organisations (provider) {
    return this.http.get(`/source/${provider}/orgs`).then(res => res.data)
  }

  /**
   * repositories, load source repositories for given organisation
   * @param {string} provider
   * @param {string} organisation
   * @returns {Promise<[{provider: string, organisation: string, name: string}]>}
   */
  repositories (provider, organisation) {
    const params = {org: organisation}
    return this.http.get(`/source/${provider}/repos`, {params}).then(res => res.data)
  }
}
