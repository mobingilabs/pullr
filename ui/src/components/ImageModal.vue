<template>
  <Modal :title="image.name" v-model="visible" v-on:on-cancel="close">
    <div class="info-entry">
      <div class="label">Source Provider:</div>
      <div class="value">{{image.repository.provider}}</div>
    </div>
    <div class="info-entry">
      <div class="label">Repository:</div>
      <div class="value">{{image.repository.owner}}/{{image.repository.name}}</div>
    </div>
    <div class="info-entry">
      <div class="label">Dockerfile Path:</div>
      <div class="value">{{image.dockerfile_path}}</div>
    </div>
    <div class="info-entry">
      <div class="label">Docker Tags:</div>
      <div class="value no-padding">
        <Table size="small" :columns="buildColumns" :data="image.tags"/>
      </div>
    </div>
    <template slot="footer">
      <Button icon="clock" @click="$router.push({name: 'image-history', params: {key: image.key}})">HISTORY</Button>
      <Button type="primary" icon="edit" @click="$router.push({name: 'edit-image', params: {key: image.key}})">EDIT</Button>
    </template>
  </Modal>
</template>

<script type="text/jsx">
export default {
  name: 'ImageModal',
  props: {
    image: Object,
    value: Boolean
  },

  data () {
    return {
      visible: this.value,
      infoColumns: [
        {title: 'Label', key: 'label'},
        {title: 'Value', key: 'value'}
      ],
      buildColumns: [
        {title: 'Type', key: 'ref_type'},
        {title: 'Test', key: 'ref_test'},
        {
          title: 'Name',
          render: (h, obj) => {
            const tagName = (obj.row.name === '' && obj.row.type === 'tag') ? obj.row.test : obj.row.name
            return <span>{tagName}</span>
          }
        }
      ]
    }
  },

  watch: {
    value (newVal, oldVal) {
      this.visible = newVal
    }
  },

  methods: {
    close () {
      this.$emit('input', false)
      this.$emit('close')
    }
  }
}
</script>

<style lang="less" scoped>
  .info-entry {
    margin-bottom: 25px;

    .label {
      font-weight: bold;
      color: #888;
    }

    .value {
      font-size: 16px;
    }
  }
</style>
