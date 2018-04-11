import Vue from 'vue'
import Vuex from 'vuex'
import * as errors from './services/errors'

Vue.use(Vuex)

function newImage () {
  return {
    name: '',
    dockerfile_path: './Dockerfile',
    repository: {
      provider: '',
      owner: '',
      name: ''
    },
    tags: [
      {ref_type: 'branch', ref_test: '', name: ''}
    ]
  }
}

const initialState = {
  builds: {
    data: [],
    loading: false,
    loadErr: null
  },
  newImage: {
    data: newImage(),
    saving: false,
    saveErr: null,
    owners: {
      data: [],
      loading: false,
      loadErr: null
    },
    repositories: {
      data: {},
      loading: false,
      loadErr: false
    }
  },
  auth: {
    authenticating: true,
    authenticateErr: null,
    profile: {
      loading: false,
      loadErr: null,
      data: {
        user: {
          username: '',
          email: ''
        },
        tokens: {}
      }
    }
  },
  oauth: {
    loginURLs: {
      github: {loading: false, loadErr: null, url: null},
      gitlab: {loading: false, loadErr: null, url: null},
      bitbucket: {loading: false, loadErr: null, url: null}
    }
  },
  images: {
    loading: false,
    loadErr: null,
    lastLoad: null,
    operations: {},
    builds: {},
    data: []
  }
}

/**
 * @param {AuthService} authService
 * @param {OAuthService} oauthService
 * @param {SourceService} sourceService
 * @param {ImagesService} imagesService
 * @returns {Store}
 */
export default (authService, oauthService, sourceService, imagesService, buildService) => new Vuex.Store({
  state: {...initialState},
  getters: {
    getImageByKey: (state) => (key) => state.images.data.find(img => img.key === key),
    getImageBuild: (state) => (key) => state.images.builds[key],
    filteredImages: (state) => (query) => {
      if (!query) {
        return []
      }

      const regexp = new RegExp(query)
      return state.images.data.filter(img => regexp.test(img.name))
    }
  },
  mutations: {
    authenticateRequest (state) {
      state.auth.authenticating = true
    },
    authenticateSuccess (state) {
      state.auth.authenticating = false
      state.auth.authenticateErr = null
    },
    authenticateFailure (state, err) {
      state.auth.authenticateErr = err
      state.auth.authenticating = false
    },
    logout (state) {
      state.auth = {...initialState.auth}
    },
    loadImagesRequest (state) {
      state.images.loading = true
    },
    loadImagesSuccess (state, {images, builds}) {
      state.images.loadErr = null
      state.images.data = images
      state.images.builds = builds
      state.images.loading = false
      state.images.lastLoad = new Date()
    },
    loadImagesFailure (state, err) {
      state.images.loadErr = err
    },
    resetNewImage (state) {
      state.newImage.data = newImage()
      state.newImage.owners = []
      state.newImage.repositories = {}
    },
    updateNewImage (state, image) {
      state.newImage.data = image
    },
    saveNewImageRequest (state) {
      state.newImage.saving = true
    },
    saveNewImageSuccess (state) {
      state.newImage.data = newImage()
      state.newImage.saving = false
    },
    saveNewImageFailure (state, err) {
      state.newImage.saving = false
      state.newImage.saveErr = err
    },
    loadNewImageOwnersRequest (state) {
      state.newImage.owners.loading = true
    },
    loadNewImageOwnersSuccess (state, owners) {
      state.newImage.owners.data = owners
      state.newImage.data.repository.owner = owners[0]
      state.newImage.owners.loading = false
    },
    loadNewImageOwnersFailure (state, err) {
      state.newImage.owners.loadErr = err
      state.newImage.owners.loading = false
    },
    loadNewImageRepositoriesRequest (state) {
      state.newImage.repositories.loading = true
    },
    loadNewImageRepositoriesSuccess (state, {owner, repositories}) {
      Vue.set(state.newImage.repositories.data, owner, repositories)
      state.newImage.repositories.loading = false
    },
    loadNewImageRepositoriesFailure (state, err) {
      state.newImage.repositories.loadErr = err
      state.newImage.repositories.loading = false
    },
    getOAuthLoginURLRequest (state, provider) {
      state.oauth.loginURLs[provider].loading = true
    },
    getOAuthLoginURLSuccess (state, {provider, loginURL}) {
      state.oauth.loginURLs[provider].url = loginURL
      state.oauth.loginURLs[provider].loading = false
    },
    getOAuthLoginURLFailure (state, {provider, err}) {
      state.oauth.loginURLs[provider].loadErr = err
      state.oauth.loginURLs[provider].loading = false
    },
    loadProfileRequest (state) {
      state.auth.profile.loading = true
    },
    loadProfileSuccess (state, {user, tokens}) {
      state.auth.profile.data.user = user

      state.auth.profile.data.tokens = {}
      tokens = tokens || []
      tokens.forEach(token => {
        state.auth.profile.data.tokens[token] = true
      })

      state.auth.profile.loading = false
    },
    loadProfileFailure (state, err) {
      state.auth.profile.loadErr = err
      state.auth.profile.loading = false
    },
    deleteImageRequest (state, key) {
      state.images.operations = {
        ...state.images.operations,
        [key]: {...state.images.operations[key], deleting: true, deleteErr: null}
      }
    },
    deleteImageSuccess (state, key) {
      state.images.operations = {
        ...state.images.operations,
        [key]: {...state.images.operations[key], deleting: false, deleteErr: null}
      }

      state.images.data = state.images.data.filter(img => img.key !== key)
    },
    deleteImageFailure (state, key, err) {
      state.images.operations = {
        ...state.images.operations,
        [key]: {...state.images.operations[key], deleting: false, deleteErr: err}
      }
    },
    updateImage (state, {key, data}) {
      state.images.data = state.images.data.map(img => {
        if (img.key === key) {
          return data
        }

        return img
      })
    },
    loadBuildsRequest (state) {
      state.builds.loading = true
    },
    loadBuildsSuccess (state, builds) {
      state.builds.data = builds
      state.builds.loading = false
    },
    loadBuildsFailure (state, err) {
      state.builds.loadErr = err
      state.builds.loading = false
    }
  },
  actions: {
    async authenticate ({commit, dispatch}) {
      commit('authenticateRequest')
      try {
        await authService.authenticate()
        await dispatch('loadProfile')
        commit('authenticateSuccess')
      } catch (err) {
        commit('authenticateFailure', err)
        throw err
      }
    },
    async login ({commit}, {username, password}) {
      commit('authenticateRequest')
      try {
        await authService.login(username, password)
        commit('authenticateSuccess')
      } catch (err) {
        commit('authenticateFailure', err)
        throw err
      }
    },
    async logout ({commit}) {
      commit('logout')
      authService.logout()
      return Promise.resolve()
    },
    async fetchImages ({commit}) {
      commit('loadImagesRequest')
      try {
        const data = await imagesService.list()
        commit('loadImagesSuccess', data)
      } catch (err) {
        commit('loadImagesFailure', err)
        throw err
      }
    },
    async updateNewImage ({commit, state, dispatch}, image) {
      const previousOwner = state.newImage.data.repository.owner
      const previousProvider = state.newImage.data.repository.provider

      commit('updateNewImage', image)

      if (previousProvider !== image.repository.provider) {
        await dispatch('loadNewImageOwners')
      }
      if (previousOwner !== image.repository.owner) {
        await dispatch('loadNewImageRepositories', image.repository.owner)
      }
    },
    async saveNewImage ({commit, state}) {
      commit('saveNewImageRequest')
      try {
        const image = await imagesService.create(state.newImage.data)
        commit('saveNewImageSuccess', image)
      } catch (err) {
        commit('saveNewImageFailure', err)
        throw err
      }
    },
    async loadNewImageOwners ({commit, state}) {
      commit('loadNewImageOwnersRequest')
      try {
        const owners = await sourceService.owners(state.newImage.data.repository.provider)
        commit('loadNewImageOwnersSuccess', owners)
      } catch (err) {
        commit('loadNewImageOwnersFailure', err)
        throw err
      }
    },
    async loadNewImageRepositories ({commit, state}, owner) {
      commit('loadNewImageRepositoriesRequest')
      try {
        const repositories = await sourceService.repositories(state.newImage.data.repository.provider, owner)
        commit('loadNewImageRepositoriesSuccess', {owner, repositories})
      } catch (err) {
        commit('loadNewImageRepositoriesFailure', err)
        throw err
      }
    },
    async loadImage ({commit, state, dispatch, getters}, key) {
      let image = getters.getImageByKey(key)
      if (image) {
        return Promise.resolve(image)
      }

      await dispatch('fetchImages')
      image = getters.getImageByKey(key)
      if (!image) {
        throw new errors.ApiError(errors.KindNotFound, 'Image not found')
      }

      return Promise.resolve(image)
    },
    async loadBuilds ({commit}) {
      commit('loadBuildsRequest')
      try {
        const builds = await buildService.list()
        commit('loadBuildsSuccess', builds)
      } catch (err) {
        commit('loadBuildsFailure', err)
        throw err
      }
    },
    async loadImageBuilds (context, {key, page, perPage}) {
      return buildService.history(key, {page, per_page: perPage})
    },
    async updateImage ({commit, state}, {key, data}) {
      const updatedImage = await imagesService.update(key, data)
      commit('updateImage', {key, data: updatedImage})
    },
    async getOAuthLoginURL ({commit}, provider) {
      commit('getOAuthLoginURLRequest', provider)
      try {
        const loginURL = await oauthService.loginURL(provider)
        commit('getOAuthLoginURLSuccess', {provider, loginURL})
      } catch (err) {
        commit('getOAuthLoginURLFailure', {provider, err})
        throw err
      }
    },
    async loadProfile ({commit}) {
      commit('loadProfileRequest')
      try {
        const profile = await authService.profile()
        commit('loadProfileSuccess', profile)
      } catch (err) {
        commit('loadProfileFailure', err)
        throw err
      }
    },
    async deleteImage ({commit}, key) {
      commit('deleteImageRequest', key)
      try {
        await imagesService.delete(key)
        commit('deleteImageSuccess', key)
      } catch (err) {
        commit('deleteImageFailure', key)
      }
    }
  }
})
