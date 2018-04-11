<template>
  <div class="login">
    <h1 class="title">Login</h1>
    <i-form @submit="login" ref="form" :model="form" :rules="loginRules">
      <FormItem label="Username" prop="username">
        <Input v-model="form.username" size="large" icon="person"/>
      </FormItem>
      <FormItem class="password" label="Password" prop="password">
        <template slot="label">
          Password
          <router-link :to="{name: 'forgot-password'}" class="forgot-password">Forgot your password?</router-link>
        </template>
        <Input v-model="form.password" size="large" type="password" icon="key"/>
      </FormItem>
    </i-form>
    <div class="actions">
      <Button :loading="inProgress" type="primary" @click="login">LOGIN</Button>
      <div v-if="serverError" class="server-error">
        <Icon type="alert-circled"/>
        {{serverError}}
      </div>
    </div>
    <div class="divider"></div>
    <p class="invite">Don't have an account yet?</p>
    <Button @click="$router.push({name: 'register'})" type="ghost">REGISTER</Button>
  </div>
</template>

<script>
import * as errors from '../services/errors'

export default {
  name: 'Login',
  inject: ['authService'],
  data () {
    return {
      form: {
        username: '',
        password: ''
      },
      inProgress: false,
      serverError: null,
      loginRules: {
        username: [{required: true, message: 'Please enter a username', trigger: 'blur'}],
        password: [{required: true, message: 'Please enter a password', trigger: 'blur'}]
      }
    }
  },
  methods: {
    async login () {
      const valid = await this.$refs.form.validate()
      if (!valid) {
        return this.$Notice.error({
          title: 'Error',
          desc: 'Please fix the errors mentioned'
        })
      }

      this.inProgress = true
      this.serverError = null
      try {
        await this.$store.dispatch('login', {username: this.form.username, password: this.form.password})
        this.$router.push({name: 'images'})
      } catch (err) {
        this._onError(err)
      } finally {
        this.inProgress = false
      }
    },
    _onError (err) {
      let desc = 'Some unexpected error happened.'
      if (err.kind === errors.KindUnauthorized) {
        desc = 'Username or password is invalid'
      }

      this.serverError = desc
    }
  }
}
</script>

<style scoped>
  .forgot-password {
    float: right;
    clear: both;
  }

  .password >>> label {
    display: block;
    text-align: left;
    width: 100%;
    padding-right: 0;
  }

  .password >>> .ivu-form-item-content {
    display: block;
    clear: both;
  }
</style>
