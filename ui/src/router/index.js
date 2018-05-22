import Vue from 'vue'
import Router from 'vue-router'

import Images from '../pages/Images'
import CreateImage from '../pages/create-image/CreateImage'
import EditImage from '../pages/EditImage'
import NotFound from '../pages/NotFound'
import Login from '../pages/Login'
import Register from '../pages/Register'
import OAuthRedirect from '../pages/OAuthRedirect'
import History from '../pages/History'
import ImageHistory from '../pages/ImageHistory'

Vue.use(Router)

/**
 * @param {AuthService} authService
 * @returns {VueRouter}
 */
export default (authService) => {
  const router = new Router({
    mode: 'history',
    routes: [
      {
        path: '/images/create',
        name: 'create-image',
        component: CreateImage,
        meta: {group: 'images', title: 'Create Image', requiresAuth: true}
      },
      {
        path: '/images/:key?',
        alias: '/',
        name: 'images',
        component: Images,
        meta: {group: 'images', title: 'Images', requiresAuth: true}
      },
      {
        path: '/images/:key/edit',
        name: 'edit-image',
        component: EditImage,
        meta: {group: 'images', title: 'Edit Image', requiresAuth: true}
      },
      {
        path: '/history',
        name: 'history',
        component: History,
        meta: {group: 'history', title: 'History', requiresAuth: true}
      },
      {
        path: '/history/:key',
        name: 'image-history',
        component: ImageHistory,
        meta: {group: 'history', title: 'History', requiresAuth: true}
      },
      {
        path: '/auth/login',
        name: 'login',
        component: Login,
        meta: {group: 'auth'}
      },
      {
        path: '/auth/register',
        name: 'register',
        component: Register,
        meta: {group: 'auth'}
      },
      {
        path: '/oauth/cb',
        name: 'oauth',
        component: OAuthRedirect,
        meta: {group: 'auth'}
      },
      {
        path: '*',
        name: 'not-found',
        component: NotFound
      }
    ]
  })

  router.beforeEach((to, from, next) => {
    authService.ready.catch(() => {}).then(() => {
      if (to.name === 'login' && authService.authenticated) {
        return next({name: 'images'})
      }

      if (to.matched.some(record => record.meta.requiresAuth)) {
        if (authService.authenticated) {
          return next()
        }

        return next({name: 'login'})
      }

      return next()
    })
  })

  return router
}
