<template>
  <div>
    <Table class="images-table" size="large" :columns="columns" :data="images" no-data-text="There are no images yet"/>
    <ImageModal v-model="showImageModal" :image="activeImage" @close="onImageModalClosed"/>
    <Spin fix v-if="imagesLoading"/>
  </div>
</template>

<script type="text/jsx">
import Vue from 'vue'
import { Button, Tooltip } from 'iview'
import { mapState, mapGetters, mapActions } from 'vuex'
import ImageModal from '../components/ImageModal'

export default {
  name: 'Images',
  components: {
    ImageModal
  },
  computed: {
    ...mapState({
      images: (state) => state.images.data || [],
      imagesLoading: (state) => state.images.loading,
      loadErr: (state) => state.images.loadErr,
      imageOperations: (state) => state.images.operations
    }),
    ...mapGetters(['getImageBuild'])
  },
  watch: {
    loadErr (newVal) {
      if (!newVal) {
        return
      }

      this.$Notice.error({title: 'Failed to load images', desc: newVal.message})
    },
    '$route.params.key' (key) {
      if (key) {
        this.openImage(key)
      } else {
        this.showImageModal = false
      }
    }
  },
  data () {
    return {
      showImageModal: false,
      activeImage: {repository: {}, tags: []},
      deleteOperations: {},
      columns: [
        {
          title: 'Name',
          render: (h, obj) => {
            const name = obj.row.name
            const key = obj.row.key
            return <a onClick={() => this.openImage(key)}>{name}</a>
          }
        },
        {
          title: 'Repository',
          render: (h, obj) => {
            const repo = obj.row.repository
            return <span>{repo.provider}/{repo.owner}/{repo.name}</span>
          }
        },
        {
          title: 'Tags',
          render: (h, obj) => {
            const tags = obj.row.tags.map(tag => {
              if (tag.name === '' && tag.ref_type === 'tag') {
                return tag.ref_test
              }
              return tag.name
            }).join(', ')

            return <span>{tags}</span>
          }
        },
        {
          title: 'Last Build',
          className: 'last-build',
          render: (h, obj) => {
            const build = this.getImageBuild(obj.row.key)
            if (!build) {
              return null
            }

            const time = build.status === 'in_progress' ? build.started_at : build.finished_at
            return <span>{Vue.filter('date')(time)}</span>
          }
        },
        {
          title: ' ',
          align: 'right',
          render: (h, obj) => {
            const editRoute = {name: 'edit-image', params: {key: obj.row.key}}
            const historyRoute = {name: 'image-history', params: {key: obj.row.key}}
            const operations = this.imageOperations[obj.row.key] || {}
            return (
              <div class="row-actions">
                <Tooltip content="Delete">
                  <Button size="small" loading={operations.deleting} onClick={() => this.deleteImage(obj.row.key)}
                    icon="trash-b"/>
                </Tooltip>
                <Tooltip content="History">
                  <Button size="small" onClick={() => this.$router.push(historyRoute)} icon="clock"/>
                </Tooltip>
                <Tooltip content="Edit">
                  <Button type="primary" size="small" onClick={() => this.$router.push(editRoute)} icon="edit"/>
                </Tooltip>
              </div>
            )
          }
        }
      ]
    }
  },
  methods: {
    ...mapActions(['deleteImage']),
    openImage (key) {
      const image = this.$store.getters.getImageByKey(key)
      if (!image) {
        return
      }

      this.activeImage = image
      this.showImageModal = true
      this.$router.push({name: 'images', params: {'key': key}})
    },
    onImageModalClosed () {
      this.$router.push({name: 'images'})
    }
  },
  async mounted () {
    await this.$store.dispatch('fetchImages')
    const key = this.$route.params['key']
    if (key) {
      this.openImage(key)
    }
  }
}
</script>

<style lang="less">
  .images-table {
    .row-actions {
      display: none;
    }

    tr:hover .row-actions {
      display: block;

      .ivu-btn {
        margin: 0 5px;
        width: 32px;
        height: 32px;
      }
    }
  }
</style>
