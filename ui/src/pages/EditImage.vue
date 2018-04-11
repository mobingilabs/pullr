<template>
  <div style="min-height: 300px">
    <Card dis-hover>
      <p slot="title"><span v-if="!!image">Edit {{image.name}}</span></p>
      <ImageEditor v-if="!!image" :image.sync="image" @cancel="cancel()" @save="save()"/>
      <Spin fix v-if="!image || saveInProgress"/>
    </Card>
  </div>

</template>

<script>
import * as errors from '../services/errors'
import ImageEditor from '../components/ImageEditor'

export default {
  name: 'EditImage',
  components: {
    ImageEditor
  },
  data () {
    return {
      image: null,
      notFound: false,
      saveInProgress: false
    }
  },
  methods: {
    cancel () {
      this.$router.push({name: 'images'})
    },
    async save () {
      this.saveInProgress = true
      try {
        await this.$store.dispatch('updateImage', {key: this.$route.params.key, data: this.image})
        this.saveInProgress = false
        this.$Message.success('Image successfully updated')
      } catch (err) {
        this.$Notice.error({title: 'Failed to update image', desc: err.message})
        this.saveInProgress = false
      }
    }
  },
  async created () {
    try {
      await this.$store.dispatch('loadImage', this.$route.params.key)
      const image = this.$store.getters.getImageByKey(this.$route.params.key)
      this.image = {
        ...image,
        repository: {...image.repository},
        tags: [...image.tags.map(tag => ({...tag}))]
      }
    } catch (err) {
      if (err.kind === errors.KindNotFound) {
        this.notFound = true
      } else {
        this.$Notice.error({title: 'Something bad happened', desc: err.message})
      }
    }
  }
}
</script>

<style scoped>

</style>
