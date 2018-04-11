import { EventEmitter } from 'events'
import { ApiError, KindUnauthorized } from './errors'

const XAuthTokenHeader = 'x-auth-token'
const XRefreshTokenHeader = 'x-refresh-token'

export default class AuthService extends EventEmitter {
  /**
   * @param {AxiosInstance} http Http library, should implement axios interface
   * @param {string} authToken Initial authToken, useful for setting it from localStorage
   * @param {string} refreshToken Initial refreshToken useful for setting it from localStorage
   */
  constructor (http, authToken, refreshToken) {
    super()
    this.http = http
    this.authToken = authToken
    this.refreshToken = refreshToken
    this.authenticated = false
    this.ready = null

    this.http.interceptors.request.use(this._interceptRequest)
    this.http.interceptors.response.use(this._interceptResponse, this._interceptResponseError)
  }

  /**
   * authenticate authenticates the user by calling profile with stored authToken and refreshToken
   * @returns {Promise<null>}
   */
  authenticate = () => {
    this.ready = this.profile().then(() => {
      this.authenticated = true
    }).catch(err => {
      this.authenticated = false
      this.authToken = null
      this.refreshToken = null
      throw err
    })

    return this.ready
  }

  /**
   * login, logs in a user with given username and password
   * @param {string} username
   * @param {string} password
   * @returns {Promise<{user: Object, tokens: string[]}>}
   */
  login = async (username, password) => {
    await this.http.post('/auth/login', {username, password})
    this.authenticated = true
    this.ready = Promise.resolve()
    this.emit('login')
  }

  /**
   * register, registers a new user
   * @param {string} username
   * @param {string} email
   * @param {string} password
   * @returns Promise<null>
   */
  register = (username, email, password) => {
    return this.http.post('/auth/register', {username, email, password}).then(() => null)
  }

  /**
   * logout, clears authentication data and logs out the user
   */
  logout = () => {
    this.authToken = null
    this.refreshToken = null
    this.authenticated = false
    this.emit('logout')
  }

  /**
   * profile fetches the authenticated user profile.
   * @returns {Promise<{user: Object, tokens: string[]}>}
   */
  profile = () => {
    return this.http.get('/user/profile').then(res => res.data)
  }

  _interceptRequest = (config) => {
    if (!this.authToken && !this.refreshToken) {
      return config
    }

    config.headers['Authorization'] = `Bearer ${this.authToken}`
    config.headers[XRefreshTokenHeader] = this.refreshToken
    return config
  }

  _interceptResponse = (res) => {
    if (res.headers[XAuthTokenHeader] && res.headers[XRefreshTokenHeader]) {
      this.authToken = res.headers[XAuthTokenHeader]
      this.refreshToken = res.headers[XRefreshTokenHeader]
      this.emit('tokens', this.authToken, this.refreshToken)
    }

    return res
  }

  _interceptResponseError = (err) => {
    if (!err.response) {
      return Promise.reject(err)
    }

    if (err.response.data.kind === KindUnauthorized) {
      this.logout()
    }

    return Promise.reject(new ApiError(err.response.data.kind, err.response.data.message))
  }
}
