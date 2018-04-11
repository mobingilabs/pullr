import EventEmitter from 'events'

export default class OAuthService extends EventEmitter {
  /**
   * @param {AxiosInstance} http
   */
  constructor (http) {
    super()
    this.http = http

    if (window) {
      window.addEventListener('message', this._onWindowMessage, false)
    }
  }

  /**
   * loginURL fetches oauth login url for the given oauth provider
   * @param {string} provider
   * @returns {Promise<string>}
   */
  loginURL (provider) {
    const params = {redirect: `${location.protocol}//${location.host}/oauth/cb`}
    return this.http.get(`/oauth/${provider}/login_url`, {params}).then(res => res.data.login_url)
  }

  _onWindowMessage = (e) => {
    if (e.origin !== location.origin) {
      return
    }

    if (e.data === 'OAUTH_SUCCESS') {
      this.emit('login')
    }
  }
}
