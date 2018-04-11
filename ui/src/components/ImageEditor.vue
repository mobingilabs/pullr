<template>
  <div>
    <Form ref="imageForm" :model="image" :rules="imageRules">
      <Row :gutter="16">
        <Col :sm="24" :md="12">
          <Row>
            <Col :sm="24" :md="12">
              <FormItem label="SOURCE">
                <Input readonly :value="source"/>
              </FormItem>
            </Col>
            <Col :sm="24" :md="12">
              <FormItem prop="name" label="IMAGE NAME">
                <Input @input="onImageNameUpdate($event)" :value="image.name"/>
              </FormItem>
            </Col>
          </Row>
          <Row>
            <Col :sm="24">
              <FormItem prop="dockerfile_path" label="DOCKERFILE PATH">
                <Input @input="onDockerfilePathUpdate($event)" :value="image.dockerfile_path"/>
              </FormItem>
            </Col>
          </Row>
        </Col>
        <Col class="docker-tags" :sm="24" :md="12">
          <FormItem prop="tags" label="DOCKER TAGS">
            <table class="tags-table" cellpadding="0" cellspacing="0">
              <thead>
              <tr>
                <th>TYPE</th>
                <th>TEST</th>
                <th>NAME</th>
                <th></th>
              </tr>
              </thead>
              <tbody>
              <tr v-for="(tag, i) in image.tags" :key="i">
                <td>
                  <Select :value="tag.ref_type" @input="onDockerTagUpdate(i, 'ref_type', $event)">
                    <Option value="branch">Branch</Option>
                    <Option value="tag">Tag</Option>
                  </Select>
                </td>
                <td>
                  <FormItem :prop="`tags.${i}.ref_test`" :rules="[{required: true, message: 'Can not be empty'}]">
                    <Input :value="tag.ref_test"
                           :placeholder="tag.ref_type === 'branch' ? 'master' : '/.*/'"
                           @input="onDockerTagUpdate(i, 'ref_test', $event)"
                    />
                  </FormItem>
                </td>
                <td>
                  <FormItem :prop="`tags.${i}.name`" :rules="[{validator: (rule, value, cb) => validateTagName(i, cb), required: true, message: 'Can not be empty'}]">
                    <Input :value="tag.name"
                           :placeholder="tag.ref_type === 'branch' ? 'latest' : 'same as tag'"
                           @input="onDockerTagUpdate(i, 'name', $event)"
                    />
                  </FormItem>
                </td>
                <td>
                  <Button v-if="i > 0" shape="circle" size="small" @click="onDockerTagRemove(i)" icon="close"/>
                </td>
              </tr>
              </tbody>
            </table>
            <Button type="primary" size="small" icon="plus" @click="onDockerTagAdd()">ADD TAG</Button>
          </FormItem>
        </Col>
      </Row>
    </Form>
    <Row :gutter="16" class="actions">
      <Col :xs="24" pull="right">
        <Button type="ghost" @click="$emit('cancel')">CANCEL</Button>
        <Button type="primary" :loading="saveInProgress" @click="handleSave">SAVE</Button>
      </Col>
    </Row>
  </div>
</template>

<script>
export default {
  name: 'ImageEditor',
  props: {
    image: {type: Object, required: true},
    saveInProgress: {type: Boolean}
  },
  computed: {
    source () {
      const repo = this.image.repository
      return `${repo.provider}/${repo.owner}/${repo.name}`
    }
  },
  data () {
    return {
      imageRules: {
        name: [{required: true, trigger: 'blur', message: 'An image name is required'}],
        dockerfile_path: [{required: true, trigger: 'blur', message: 'A dockerfile path is required'}]
      }
    }
  },
  mounted () {
    if (this.image.name === '') {
      return this._updateImage({...this.image, name: this.image.repository.name})
    }
  },
  methods: {
    _updateImage (image) {
      this.$emit('update:image', image)
    },
    async handleSave () {
      const valid = await this.$refs.imageForm.validate()
      if (valid) {
        this.$emit('save')
      }
    },
    validateTagName (index, callback) {
      if (this.image.tags[index].ref_type !== 'tag' && this.image.tags[index].name === '') {
        return callback(new Error('Name cannot be empty'))
      }

      return callback()
    },
    onImageNameUpdate (name) {
      this._updateImage({...this.image, name})
    },
    onDockerfilePathUpdate (dockerfilePath) {
      this._updateImage({...this.image, dockerfile_path: dockerfilePath})
    },
    onDockerTagUpdate (index, field, value) {
      const tags = this.image.tags.map((tag, i) => {
        if (i === index) {
          return {...tag, [field]: value}
        }
        return tag
      })

      this._updateImage({...this.image, tags})
    },
    onDockerTagAdd () {
      const tags = [].concat(this.image.tags, [{ref_type: 'branch', ref_test: '', name: ''}])
      this._updateImage({...this.image, tags})
    },
    onDockerTagRemove (index) {
      const tags = [...this.image.tags]
      tags.splice(index, 1)
      this._updateImage({...this.image, tags})
    }
  }
}
</script>

<style lang="less" scoped>
  .docker-tags {
    .ivu-form-item-content {
      clear: both;
    }

    .tags-table {
      width: 100%;
      border: 0;

      th, td {
        text-align: left;
        line-height: 16px;
        padding-right: 6px;
        color: #aaa;
        padding-bottom: 6px;
      }

      th:last-child {
        min-width: 40px;
      }

      tbody {
        tr {
          height: 60px;
          vertical-align: top;
        }
      }
    }
  }
</style>
