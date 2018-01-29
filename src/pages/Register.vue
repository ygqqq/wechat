<template>
  <div>
    <img src="../assets/logo.png">
    <mt-field label="用户名" placeholder="请输入用户名" v-model="username"></mt-field>
    <mt-field label="收信人" placeholder="请输入收信人" v-model="rec_user"></mt-field>

    <mt-field label="密码" placeholder="请输入密码" type="password" v-model="password"></mt-field>
    <mt-field label="密码确认" placeholder="请再次输入密码" type="password" v-model="passwordConfirm"></mt-field>
    <mt-field label="昵称" placeholder="请输入昵称" v-model="nickname"></mt-field>



    <mt-radio class="sex-radio"
        title="性别"
        v-model="sex"
        :options="[{label: '男',value: '0'},{label: '女',value: '1'}]">
    </mt-radio>
    <div  class="oprate-area">

        <mt-button type="primary" size="large" @click="submit">提交注册</mt-button>
        <br>
         <mt-button type="primary" size="large" @click="send">添加好友</mt-button>
        <br>
         <mt-button type="primary" size="large" @click="init">初始化ws</mt-button>
    </div>
    
  </div>
</template>

<script>
import config from "../../config/local.config"
import axios from "axios"
export default {
  data(){
      return{
          username: '',
          password: '',
          passwordConfirm: '',
          nickname: '',
          sex: '0',
          rec_user: ''
          
      }
  },
  methods:{
      submit(){
          
          axios.post('/api/user/register', {
                username: this.username,
                password: this.password,
                gender:this.sex,
                nickname:this.nickname
            })
            .then(function (response) {
                let res = response.data;
                //if(res.success){
                    alert(res.msg)
                //}
            })
            .catch(function (error) {
                console.log(error);
            });
      },
      init(){
        this.ws = new WebSocket(config.wsUrl+'?a='+this.username);
        this.ws.addEventListener('message', function(e) {
            var msg = JSON.parse(e.data);
            console.log(msg)
        });
      },
      send(){
          
        //console.log(this.ws)
 
        this.ws.send(
            JSON.stringify({
                src: this.username,
                dst: this.rec_user,
                message: this.password,
                messagetype: 4
            }
        ));
      },
      getCookie(name) {
        var arr, reg = new RegExp("(^| )" + name + "=([^;]*)(;|$)");
        if (arr = document.cookie.match(reg))
            return (arr[2]);
        else
            return null;
    }
  },
  created: function() {
      var self = this;
      //let name = this.getCookie("username")
      //this.ws = new WebSocket('ws://127.0.0.1:8000/user/ws?a='+name);
  }
}
</script>
<style>
    .sex-radio{
        display: flex;
        justify-content: flex-start;
        padding-left: 50px;
    }
    .mint-radiolist-title{
        font-size: 16px;
        color: #333;
    }
    .oprate-area{
        margin-top: 35px;
        text-align: center;
    }
</style>
