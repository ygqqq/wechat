<template>
    <div>
      <mt-header fixed :title="title">
        <router-link to="/" slot="left">
          <mt-button icon="back">返回</mt-button>
        </router-link>
        <!-- <mt-button icon="more" slot="right" class="hasMoreButton">
        </mt-button> -->
        <mt-button slot="right" class="addFriend" @click="addFriend">
          +
        </mt-button>
      </mt-header>
      <div class="message-box">
        <ul>
            <li v-for="item in myFriends" :key="item.Id">
              <span>昵称：{{item.NickName}}</span>
              <span>登陆名：{{item.UserName}}</span>
            </li>
          </ul>
      </div>
    </div>
  </template>
  
  <script>
  import Vue from "vue"
  import axios from "axios"
  import config from "../../config/local.config"
  import { Header } from 'mint-ui'
  import { MessageBox } from 'mint-ui'

  Vue.component(Header.name, Header)

  const ErrorMsg = 0  //错误消息
  const OnlineRemind	= 1	//上线提醒
  const OfflineRemind   = 2 //下线提醒 
  const AddFriendReq	= 3 //添加好友请求
  const AgreeAdd		= 4 //同意好友请求
  const DisAgreeAdd 	= 5 //拒绝好友请求
  const NormalMsg		= 10 //普通通知消息

  export default{
    name: 'message',
    data () {
      return {
        nowTitle: '消息',
        username: this.getCookie("username"),
        myFriends: [],
        getData: () => {
          const _this = this
          axios.get('/api/user/friends?name='+ this.username, {})
          .then(function (response) {
            if (response.data.success) {
              let friendsArr = JSON.parse(response.data.msg)
              for (var i = 0; i < friendsArr.length; i++) {
                let friends = {
                  UserName: friendsArr[i].UserName,
                  NickName: friendsArr[i].NickName,
                  Id: friendsArr[i].Id_,
                  Status: friendsArr[i].Status,
                  CreateAt: friendsArr[i].CreateAt,
                }
                _this.myFriends.push(friends)
              }
            } else {
              console.log('false')
            }
          })
          .catch(function (error) {
            alert(error);
          });
        }
      }
    },
    computed:{
      title(){
        return this.getCookie("username")
      }
    },
    methods: {
      addFriend () {
        MessageBox.prompt('请输入好友姓名','').then(({ value, action }) => {
          if( action == 'confirm') {
            if (value == null || value.length <= 0) {
              MessageBox('', '好友昵称不能为空') 
              return
            }
            console.log(value,this.username)
            if (value == this.username) {
              MessageBox('', '不能添加自己') 
              return
            }
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
    created () {
      var _this = this
      let username = this.getCookie("username")

      _this.ws = new WebSocket(config.wsUrl+'?a='+username) //注册WebSocket
      this.ws.addEventListener('message', function(e) {  //监听WebSocket
        var msg = JSON.parse(e.data);
        switch(msg.messagetype){
          case 3:     
            _this.agreeFriend(msg.src)
            break
          case 10:     
            MessageBox('', msg.message)
            break
          // default:
          //   MessageBox('', msg.message)
          //   break
        }
      })

      this.getData()
    }
  }
  </script>
  
  <style lang="scss" scoped type="text/css">
    .addFriend{
      font-size: 26px;
    }
    .message-box{
      margin-top: 40px;

      ul li{
        line-height: 40px;
        height: 40px;
        border-bottom: 1px solid #999;
      }
    }
  </style>