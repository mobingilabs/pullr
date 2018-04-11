<template>
  <div class="select-provider">
    <Row :gutter="20" type="flex" class="row">
      <Col v-for="p in providers" :key="p.value">
        <Card class="provider-card">
          <p><strong>{{p.title}}</strong></p>
          <img class="logo" v-once :src="p.icon"/>
          <Button v-if="p.disabled" type="ghost" disabled>COMING SOON</Button>
          <Button v-else-if="p.hasToken" type="primary" @click="select(p)">SELECT</Button>
          <a v-else :href="loginURLs[p.value].url" target="_blank" class="ivu-btn ivu-btn-primary">LINK ACCOUNT</a>
          <Spin size="large" fix v-if="loginURLs[p.value].loading || isProfileLoading" />
        </Card>
      </Col>
    </Row>
  </div>
</template>

<script>
import GithubLogo from '../../assets/github-logo.svg'
import GitlabLogo from '../../assets/gitlab-logo.svg'
import BitbucketLogo from '../../assets/bitbucket-logo.svg'

export default {
  name: 'SelectSourceProvider',
  props: {
    next: {type: Function, required: true}
  },
  inject: ['oauthService'],
  computed: {
    providers () {
      const tokens = this.$store.state.auth.profile.data.tokens
      return [
        {value: 'github', title: 'Github', icon: GithubLogo, hasToken: !!tokens['github']},
        {value: 'gitlab', title: 'Gitlab', icon: GitlabLogo, hasToken: !!tokens['gitlab'], disabled: true},
        {value: 'bitbucket', title: 'Bitbucket', icon: BitbucketLogo, hasToken: !!tokens['bitbucket'], disabled: true}
      ]
    },
    image () {
      return this.$store.state.newImage.data
    },
    loginURLs () {
      return this.$store.state.oauth.loginURLs
    },
    isProfileLoading () {
      return this.$store.state.auth.profile.loading
    }
  },
  watch: {
    providers () {
      this.fetchLoginURLs()
    }
  },
  methods: {
    select (provider) {
      this.$emit('update:image', {...this.image, repository: {...this.image.repository, provider: provider.value}})
      this.next()
    },
    fetchLoginURLs () {
      this.providers.forEach(async (provider) => {
        if (provider.hasToken || provider.disabled) {
          return
        }

        this.$store.dispatch('getOAuthLoginURL', provider.value)
      })
    }
  },
  mounted () {
    this.fetchLoginURLs()
  },
  beforeDestroy () {
    window.removeEventListener('message', this.onWindowMessage)
  }
}
</script>

<style lang="less" scoped>
  .provider-card {
    width: 200px;
    text-align: center;

    .logo {
      width: 50%;
      margin: 20px auto;
    }

    button {
      display: block;
      width: 100%;
    }
  }
</style>
