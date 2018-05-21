<template>
  <div class="login">
    <h1 class="title">Register</h1>
    <i-form ref="form" :model="form" :rules="registerRules">
      <FormItem label="Username" prop="username">
        <Input v-model="form.username" size="large" icon="person"/>
      </FormItem>
      <FormItem label="E-Mail" prop="email">
        <Input v-model="form.email" size="large" icon="email" type="email"/>
      </FormItem>
      <FormItem label="Password" prop="password">
        <Input v-model="form.password" size="large" type="password" icon="key"/>
      </FormItem>
      <FormItem label="Password Again" prop="passwordAgain">
        <Input v-model="form.passwordAgain" size="large" type="password" icon="key"/>
      </FormItem>
    </i-form>
    <div class="actions">
      <Button :loading="inProgress" type="primary" @click="handleRegister">REGISTER</Button>
      <div v-if="serverError" class="server-error"><Icon type="alert-circled"/> {{serverError}}</div>
    </div>
    <div class="divider"></div>
    <p class="invite">Already registered?</p>
    <Button @click="$router.push({name: 'login'})" type="ghost">LOGIN</Button>
  </div>
</template>

<script>
import { mapActions } from 'vuex'
import * as errors from '../services/errors'

export default {
  name: 'Register',
  data () {
    return {
      form: {
        username: '',
        email: '',
        password: '',
        passwordAgain: ''
      },
      serverError: null,
      inProgress: false,
      registerRules: {
        username: [
          {required: true, message: 'Please enter a username', trigger: 'blur'}
        ],
        email: [
          {required: true, type: 'email', message: 'Please enter a valid email address', trigger: 'blur'}
        ],
        password: [
          {required: true, message: 'Please enter a password', trigger: 'blur'},
          {min: 6, message: 'Password should be at least 6 characters', trigger: 'blur'}
        ],
        passwordAgain: [
          {required: true, message: 'Please retype the password above', trigger: 'blur'},
          {validator: this.matchPasswords, message: `Passwords doesn't match`, trigger: 'change'}
        ]
      }
    }
  },
  methods: {
    ...mapActions(['register']),
    async handleRegister () {
      const valid = await this.$refs.form.validate()
      if (!valid) {
        return this.$Notice.error({
          title: 'Error',
          desc: 'Please fix the errors mentioned'
        })
      }

      try {
        await this.register({username: this.form.username, email: this.form.email, password: this.form.password})
        this.$router.push('/')
      } catch (err) {
        this._onError(err)
      } finally {
        this.inProgress = false
      }
    },

    _onError (err) {
      let desc = 'Some unexpected error happened.'
      if (err.kind === errors.KindConflict) {
        if (/username/.test(err.message)) {
          desc = 'Username is already registered'
        } else {
          desc = 'Email is already registered'
        }
      } else {
        console.error(err)
      }

      this.serverError = desc
    },

    matchPasswords (rule, value) {
      if (value === this.form.password) {
        return Promise.resolve()
      }

      // eslint-disable-next-line prefer-promise-reject-errors
      return Promise.reject(null)
    }
  }
}
</script>

<style scoped>
</style>
