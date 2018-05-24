<template>
  <ImageEditor :image="image" v-on:update:image="$emit('update:image', $event)" @cancel="$emit('cancel')" :save-inprogress="saveInProgress" @save="next()" />
</template>

<script type="text/jsx">
import { mapState } from 'vuex'
import ImageEditor from '../../components/ImageEditor'

export default {
  name: 'ConfigureImage',
  components: {
    ImageEditor
  },
  props: {
    next: {type: Function, required: true}
  },
  computed: mapState({
    saveInProgress: (state) => state.newImage.saving,
    image: (state) => state.newImage.data
  })
}
</script>

<style>
</style>
