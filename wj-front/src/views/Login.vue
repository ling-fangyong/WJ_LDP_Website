<!--登录网站-->
<template>
    <div class="LoginBackgroud">
        <el-form :rules="rules" ref="LoginForm" :model="LoginForm" class="LoginContainer">
            <h3 class="LoginTitle">网站登录</h3>
            <el-form-item prop="username">
                <el-input type="text" v-model="LoginForm.username" aria-placeholder="请输入用户名">
                    <i class="el-icon-user" slot="prefix"></i>
                </el-input>
            </el-form-item>
            <el-form-item prop="password">
                <el-input type="password" v-model="LoginForm.password" aria-placeholder="请输入密码" show-password>\
                    <i class="el-icon-lock" slot="prefix"></i>
                </el-input>
            </el-form-item>
            <el-form-item style="text-align: center">
                <el-button type="primary" @click="submitLogin">登录</el-button>
            </el-form-item>
            <!-- <design ref="design"></design> -->
            <!-- <el-button @click="test">测试</el-button>  测试design未定义问题 -->
            <el-link type="primary" :underline="false" href="/register">注册新账号</el-link>
        </el-form>                
    </div>    
</template>

<script>
import * as API from "../api/user"
// import Design from './design.vue'
export default {
    name:"Login",
    // components:{
    //     Design,
    // },
    data(){
        return {
            LoginForm:{
                username:'',//用户名
                password:'',//密码
            },
            rules:{
                username:[
                    {required:true,message:'用户名不能为空',trigger:'blur'},
                    {min:4,max:16,message:'用户名长度应为4-16位',trigger:'blur'},
                ],
                password:[
                    {required:true,message:'密码不能为空',trigger:'blur'},
                    {min:6,max:20,message:'密码长度应为6-20位',trigger:'blur'},
                ]
            }
        }
    },
    methods:{
        submitLogin(){
            this.$refs.LoginForm.validate((valid)=>{
                if(valid){
                    API.login(this.LoginForm).then(res=>{
                        console.log(res);
                        if(res.code==200){
                            this.$notify({
                                type :'sucess',
                                message :this.LoginForm.username +res.msg +'!',
                            });
                            sessionStorage.setItem("username",this.LoginForm.username);
                            this.$emit("state");
                            this.$router.push({path:'/home'});
                        }else{
                            this.$notify({
                                type : 'error',
                                message : res.msg + '!',
                            });
                        }
                    }).catch((error) =>{
                        this.$notify.error({
                            message : error,
                        });
                    })
                }else{
                    this.$notify.error({
                        title:'错误',
                        message:'用户名或密码不符合要求',
                    });
                    return false;
                }
            })
        },
        // test(){
        //     console.log(this.$refs.design)
        //     this.$refs.design.sayHi();
        // }
    }
}
</script>

<style >
    .LoginBackgroud{
        position: absolute;
        width: 100%;
        height: 100%;
        background: #e0dcdc;
    }
    .LoginContainer{
        border-radius: 15px;
        background-clip: padding-box;
        margin: 188px auto;
        width: 350px;
        padding: 15px 35px 15px 35px;
        background: rgb(228, 242, 244);
        border: 1px solid #eaeaea;
        box-shadow: 0 0 25px #cac6c6;
    }
    .LoginTitle{
        margin: 0px auto 40px auto;
        text-align: center;
    }
    .el-link{
        margin-left: 41%;
        margin-top: 0%;
        line-height: 20px;
        font-size:13px;
    }
</style>