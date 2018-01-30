<template>
    <div>
        <mt-header fixed :title="nowTitle">
          <router-link to="/" slot="left">
            <mt-button icon="back">返回</mt-button>
          </router-link>
          <!-- <mt-button icon="more" slot="right" class="hasMoreButton">
          </mt-button> -->
          <mt-button slot="right" class="addFriend" @click="addFriend">
            +
          </mt-button>
        </mt-header>
      消息页面
    </div>
  </template>
  
  <script>
    import Vue from "vue"
    import { Header } from 'mint-ui'
    import config from "../../config/local.config"
    import { MessageBox } from 'mint-ui'

    Vue.component(Header.name, Header)

    const OnlineRemind	= 1	//上线提醒
    const OfflineRemind   = 2 //下线提醒 
    const AddFriendReq	= 3 //添加好友请求
    const AgreeAdd		= 4 //同意好友请求
    const DisAgreeAdd 	= 5 //拒绝好友请求

    export default{
      name: 'message',
      data () {
        return {
          nowTitle: '消息',
          username: ''
        }
      },
      methods: {
        addFriend () {
          MessageBox.prompt('请输入好友姓名','').then(({ value, action }) => {
            if( action == 'confirm') {
              ( value == null || value.length <= 0) && MessageBox('', '好友昵称不能为空')
              this.ws.send(
                JSON.stringify({
                  src: this.username,
                  dst: value,
                  messagetype: AddFriendReq
                }
              ))
            }
          })
        },
        getCookie(name) {
          var arr, reg = new RegExp("(^| )" + name + "=([^;]*)(;|$)");
          if (arr = document.cookie.match(reg))
              return (arr[2]);
          else
              return null;
        },
        agreeFriend (otherName) {
          MessageBox.confirm(otherName+'想加你好友，同意吗？').then(action => {
            if(action == 'confirm'){
              this.ws.send(
                JSON.stringify({
                  src: this.username,
                  dst: otherName,
                  messagetype: AgreeAdd
                }
              ))
            }
          })
        },
      },
      mounted () {
        var _this = this
        this.username = this.getCookie("username")

        this.ws = new WebSocket(config.wsUrl+'?a='+this.username) //注册WebSocket
        this.ws.addEventListener('message', function(e) {  //监听WebSocket
          var msg = JSON.parse(e.data);
          if( msg.messagetype == 3) {
            _this.agreeFriend(msg.src)
          }
        })
      }
    }
  </script>
  
  <style lang="scss" scoped type="text/css">
    .addFriend{
      font-size: 26px;
    }
  </style>