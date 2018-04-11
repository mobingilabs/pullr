<template>
  <MCollapse :value="step">
    <MPanel v-bind:key="s.name" v-for="(s, i) in steps" :name="s.name" disabled no-icon>
      {{i+1}}. {{s.title}}
      <template v-if="s.editable" slot="extra">
        <Button type="ghost" icon="edit" @click="gotoStep(s.name)"/>
      </template>
      <template slot="content">
        <component :is="s.component" :image="image" v-on:update:image="updateImage($event)" :next="s.next" @cancel="cancel"/>
      </template>
    </MPanel>
  </MCollapse>
</template>

<script>
import SelectSourceProvider from './SelectSourceProvider'
import SelectSourceRepository from './SelectSourceRepository'
import ConfigureImage from './ConfigureImage'

const STEP_SELECT_PROVIDER = 'source-provider'
const STEP_SELECT_REPOSITORY = 'repository'
const STEP_CONFIGURE = 'configure'

export default {
  name: 'CreateImage',
  data () {
    return {
      step: STEP_SELECT_PROVIDER,
      steps: [
        {
          name: STEP_SELECT_PROVIDER,
          title: 'SELECT A SOURCE PROVIDER',
          component: SelectSourceProvider,
          next: () => this.gotoStep(STEP_SELECT_REPOSITORY),
          editable: false,
          disabled: false
        },
        {
          name: STEP_SELECT_REPOSITORY,
          title: 'SELECT A SOURCE REPOSITORY',
          component: SelectSourceRepository,
          next: () => this.gotoStep(STEP_CONFIGURE),
          editable: false,
          disabled: true
        },
        {
          name: STEP_CONFIGURE,
          title: 'CONFIGURE IMAGE',
          component: ConfigureImage,
          next: this.saveImage,
          editable: false,
          disabled: true
        }
      ]
    }
  },
  computed: {
    image () {
      return this.$store.state.newImage.data
    }
  },
  methods: {
    cancel () {
      this.$store.commit('resetNewImage')
      this.$router.push({name: 'images'})
    },
    gotoStep (step) {
      this.step = step
      for (let i = 0; i < this.steps.length; i++) {
        if (this.steps[i].name !== step) {
          this.steps[i].editable = true
          continue
        }

        this.steps[i].editable = false
        break
      }
    },
    async saveImage () {
      try {
        await this.$store.dispatch('saveNewImage')
        this.$nextTick(() => {
          this.$router.push({name: 'images'})
        })
      } catch (err) {
        this.$Notice.error({
          title: 'Failed to create image',
          desc: err.message
        })
      }
    },
    updateImage (image) {
      this.$store.dispatch('updateNewImage', image)
    }
  }
}
</script>

<style scoped>
</style>
