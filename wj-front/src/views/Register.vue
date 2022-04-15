<!-- 用户注册 -->
<template>
    <div class="RegisterBackgroud">
        <el-form :rules="rules" ref="RegisterForm" :model="RegisterForm" class="RegisterContainer">
            <h3 class="RegisterTitle">用户注册</h3>
            <el-form-item prop="username">
                <el-input type="text" v-model="RegisterForm.username" aria-placeholder="请输入用户名">
                    <i class="el-icon-user" slot="prefix"></i>
                </el-input>
            </el-form-item>
            <el-form-item prop="password">
                <el-input type="password" v-model="RegisterForm.password" aria-placeholder="请输入密码" show-password>
                    <i class="el-icon-lock" slot="prefix"></i>
                </el-input>
            </el-form-item>
            <el-form-item prop="repassword">
                <el-input type="password" v-model="RegisterForm.repassword" aria-placeholder="再次输入密码" show-password>
                    <i class="el-icon-lock" slot="prefix"></i>
                </el-input>
            </el-form-item>
            <el-form-item style="text-align: center">
                <el-button type="primary" @click="submitRegister">注册</el-button>
            </el-form-item>
            <el-link type="primary" :underline="false" href="/Login">已有帐号，登录</el-link>
        </el-form>                
    </div>    
</template>

<script>
import * as API from "../api/user"
export default {
    name:"Register",
    data(){
        var repasswordValidate = (rule, value,callback) =>{
            if(value == ''){
                callback(new Error('重复密码不为空'))
            }else if(this.RegisterForm.repassword != this.RegisterForm.password){
                callback(new Error('两次密码不一致'))
            }else{
                callback()
            }
        }
        return {
            RegisterForm:{
                username:'',//用户名
                password:'',//密码
                repassword:'',//重复密码
            },
            rules:{
                username:[
                    {required:true,message:'用户名不能为空',trigger:'blur'},
                    {min:4,max:16,message:'用户名长度应为4-16位',trigger:'blur'},
                ],
                password:[
                    {required:true,message:'密码不能为空',trigger:'blur'},
                    {min:6,max:20,message:'密码长度应为6-20位',trigger:'blur'},
                ],
                repassword:[
                    {required:true,validator:repasswordValidate,trigger:'blur'}
                ],
            }
        }
    },
    methods:{
        submitRegister(){
            this.$refs.RegisterForm.validate((valid)=>{
                if(valid){
                    API.register(this.RegisterForm).then(res=>{
                        console.log(res);
                        if(res.code==200){
                            this.$notify({
                                type :'sucess',
                                message :this.RegisterForm.username +res.msg +'!',
                            });
                            sessionStorage.setItem("username",this.RegisterForm.username);
                            this.$emit("state");//传递状态
                            this.router.push({path:'/home'});
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
        }
    }
}
</script>

<style >
    .RegisterBackgroud{
        position: absolute;
        width: 100%;
        height: 100%;
        background: #e0dcdc;
    }
    .RegisterContainer{
        border-radius: 15px;
        background-clip: padding-box;
        margin: 188px auto;
        width: 350px;
        padding: 15px 35px 15px 35px;
        background: rgb(228, 242, 244);
        border: 1px solid #eaeaea;
        box-shadow: 0 0 25px #cac6c6;
    }
    .RegisterTitle{
        margin: 0px auto 40px auto;
        text-align: center;
    }
    .el-link{
        margin-left: 38%;
        margin-top: 0%;
        line-height: 20px;
        font-size:13px;
    }
</style>