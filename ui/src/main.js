// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import iView from 'iview'
import locale from 'iview/dist/locale/en-US'
import axios from 'axios'

import '@/assets/theme.less'
import './filters'
import Mobingi from './components/mobingi'
import App from './App'
import AuthService from './services/AuthService'
import ImagesService from './services/ImagesService'
import OAuthService from './services/OAuthService'
import SourceService from './services/SourceService'
import BuildService from './services/BuildService'
import routerFactory from './router'
import storeFactory from './store'

Vue.config.productionTip = process.env.NODE_ENV === 'production'
Vue.use(iView, {locale})
Vue.use(Mobingi)

const http = axios.create({baseURL: `https://${document.location.host}/api/v1`})
const auth = new AuthService(http, localStorage.getItem('authToken'), localStorage.getItem('refreshToken'))
const images = new ImagesService(http)
const oauth = new OAuthService(http)
const source = new SourceService(http)
const build = new BuildService(http)
const store = storeFactory(auth, oauth, source, images, build)
const router = routerFactory(auth)

auth.on('tokens', (authToken, refreshToken) => {
  localStorage.setItem('authToken', authToken)
  localStorage.setItem('refreshToken', refreshToken)
})
auth.on('logout', () => {
  localStorage.removeItem('authToken')
  localStorage.removeItem('refreshToken')
})
auth.on('login', () => store.dispatch('loadProfile'))
oauth.on('login', () => store.dispatch('loadProfile'))

store.dispatch('authenticate').catch(() => {})

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  provide: {
    oauthService: oauth,
    authService: auth
  },
  render: h => <App/>
})
