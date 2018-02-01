<template>

  <div class="login">
    <div class="login-box">
      <div class="box-flex">
        <h2>用户登录</h2>
        <mt-field label="用户名" placeholder="请输入用户名" v-model="username"></mt-field>
        <mt-field label="密码" placeholder="请输入密码" type="password" v-model="password"></mt-field>
        <div class="oprate-area">
          <mt-button type="primary" size="large" @click="login">登陆</mt-button>
          <div class="register">
            没有账号? <router-link  to="/register">立即注册</router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios"
import config from "../../config/local.config"
import { mapState,mapMutations,mapActions } from "vuex"
import { getFriends } from "../content/script/getFriends"

const storage = window.localStorage
const storageNameUser = 'dataLocalStorage__User'
export default {
  data () {
    return {
      username: '',
      password: '',
      getData: () => {
        const _this = this
        axios.get('/api/user/friends/'+ this.username, {})
        .then(function (response) {
          if (response.data.success) {
            if (response.data.msg !== 'null') {
              let friendsArr = JSON.parse(response.data.msg)
              var userMessage = {
                username: friendsArr[0].Friends[0],
                friends: []
              }
              for (var i = 0; i < friendsArr.length; i++) {
                let friends = {
                  UserName: friendsArr[i].UserName,
                  NickName: friendsArr[i].NickName,
                  Id: friendsArr[i].Id_,
                  Status: friendsArr[i].Status,
                  CreateAt: friendsArr[i].CreateAt,
                }
                userMessage.friends.push(friends)
              }
              //保存到vuex
              _this.$store.state.userMessage = userMessage
              _this.$router.push({ name: 'message'});
            }
          } else {
            console.log('false')
          }
        })
        .catch(function (error) {
          alert(error);
        })
      }
    }
  },
  methods:{
    login(){
      const _this = this
      axios.post('/api/user/login', {
          username: this.username,
          password: this.password
      })
      .then(function (response) {
        if( response.data && response.data.success) {
          //getFriends(_this)
          
          _this.getData()
          //_this.$router.push({ name: 'message'});
        } else {
          alert(response.data.msg)
        }
      })
      .catch(function (error) {
        alert(error);
      });
    }
  },
  created () {
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss" scoped type="text/css">

.register{
  color: #666;
  font-size: 14px;
  margin-top: 50px;
  display: block;

  a{
    color: #26a2ff;
  }
}
.oprate-area{
  margin-top: 35px;
  text-align: center;
}

.login{
  height: 100%;
  width: 100%;
  position: relative;
  overflow: hidden;

  .login-box{
    /* background: url(../content/images/bgLogin.jpg) no-repeat 0 0 /100%; */
    position: absolute;
    z-index: 3;
    width: 100%;
    height: 100%;
    display: flex;

    .box-flex{
      width: 88%;
      margin: auto;
      border-radius: 10px;
      padding: 20px 10px;
      background: rgba(255, 255, 255, 0.4);

      h2{
        font-size: 26px;
        font-weight: 700;
        margin-bottom: 30px;
      }
      
      .mint-cell{
        background: transparent;
        border-bottom: 1px solid #eee;
      }
    }
  }
}

.mint-field-core{ color: red}

</style>
