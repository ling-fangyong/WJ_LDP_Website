<!-- 问卷主界面 -->
<template>
    <div class="home">
        <el-row>
            <el-col :span="6">
                <div class="OperaMenu">
                    <el-tooltip class="item" effect="dark" content="创建问卷" placement="top">
                        <el-button type="text" icon="el-icon-plus" @click="addWj"></el-button>
                    </el-tooltip>
                    <el-tooltip class="item" effect="dark" content="修改问卷" placement="top">
                        <el-button type="text" icon="el-icon-edit" @click="UpdateWj"></el-button>
                    </el-tooltip>
                    <el-tooltip class="item" effect="dark" content="删除问卷" placement="top">
                        <el-button type="text" icon="el-icon-delete" @click="DeleteWj"></el-button>
                    </el-tooltip>
                    <el-tooltip class="item" effect="dark" content="分享问卷" placement="top">
                        <el-button type="text" icon="el-icon-share" @click="ShareWj"></el-button>
                    </el-tooltip>
                </div>


                <el-menu default-active="ActiveOp.toString()"  v-loading="loading" class="navbar">
                    <div style="width: 100%;text-align: center;font-size: 15px;line-height: 20px;margin-top: 20px; color: #303133" v-if="NowSelect.wjid == 0 || NowSelect.wjid == null">
                        点击创建问卷创建第一个问卷
                    </div>
                    <el-menu-item v-for="(item,index) in WjList" :key="(index+1).toString()" @click="OpenQues(index)">
                        <i class="el-icon-tickets"></i>
                        <span slot="title" style="display : inline-block">
                            {{item.title}}
                        </span>
                    </el-menu-item>
                </el-menu>
            </el-col>

            <el-col :span="18">
                <el-tabs type="border-card" v-model="ActiveTab">
                    <el-tab-pane label="问卷设计" name="wjsj">
                        <div class="WjContent">
                            <div v-show="NowSelect.wjid == 0 || NowSelect.wjid == null">请先选择问卷</div>
                            <design ref="design" v-show="NowSelect.wjid != 0 && NowSelect.wjid !=null"></design>
                        </div>
                    </el-tab-pane>
                </el-tabs>
            </el-col>

            <el-dialog title="创建问卷" :visible.sync="AddQuesShow" :close-on-click-modal="false" class="AddQues">
                <el-form ref="form" :model="AddWjmodel" label-width="80px">
                    <el-form-item label="问卷标题" style="width: 100%;" required>
                        <el-input v-model="AddWjmodel.title" placeholder="请输入问卷标题"></el-input>
                    </el-form-item>
                    <el-form-item label="问卷描述" style="width: 100%;">
                        <el-input v-model="AddWjmodel.desc" type="textarea" placeholder="请输入问卷描述" row="5"></el-input>
                    </el-form-item>
                </el-form>
                <div style="width : 100%;text-align : right">
                    <el-button style = "margin-left : 10px;" @click="AddQuesShow=false">取消</el-button>
                    <el-button type = "primary" style="margin-left : 10px;" @click = "addWjConform">确定</el-button>
                </div>
            </el-dialog>
        </el-row>
    </div>
</template>

<script>

import  * as API from '../api/question';
import Design from './design';
export default ({
    components:{
        Design,
    },
    data(){
        return{
            ActiveOp:1,
            WjList:[],
            ActiveTab:'wjsj',//问卷统计页
            AddQuesShow:false,//添加问卷弹窗
            ShareQuesShow:false,//分享问卷展示
            loading:false,//加入延迟，避免渲染为完成导致错误显示

            AddWjmodel:{
                wjid:0,
                title:'',
                desc:'感谢参加问卷调查，我们将竭力保护您的相关隐私，并按相关法律法规合理使用该问卷内容',
            },
            DelArg:{//删除问卷参数
                WjId:0,
                QuesId:0
            }
        }
    },

    computed:{
        NowSelect(){
            //console.log(this.ActiveOp)
            let now = this.WjList[this.ActiveOp -1]
            if(this.WjList==null || this.WjList.length == 0){
                return {
                    wjid:null,
                    title:null,
                    desc:null
                }
            }
            //console.log(now)
            return {
                wjid : now.wjid,
                title : now.title,
                desc : now.desc
            }
        }
    },

    mounted() {
        this.loginCheck();
    },
    methods: {
        addWj(){
            this.AddQuesShow=true;
            this.AddWjmodel = {
                title:'',
                desc:'感谢参加问卷调查，我们将竭力保护您的相关隐私，并按相关法律法规合理使用该问卷内容'
            };
        },
        addWjConform(){
            if(this.AddWjmodel.title == ''){
                this.$message({
                    id:0,
                    type:"error",
                    message:"标题不能为空",
                })
            }
            console.log(this.AddWjmodel.title)
            console.log(this.AddWjmodel.desc)        
            API.UpdateQuestionaire(this.AddWjmodel).then(res=>{
                if(res.code==200){
                    this.$message({
                        type:"success",
                        message:res.msg+'!',
                    })
                    this.getWjList();
                }else{
                    this.$message({
                        type:"error",
                        message:res.msg+'!',
                    })
                }
            }).catch(error =>{
                this.$notify.error({
                            message : error,
                });
            });
            this.AddQuesShow = false;
            this.AddWjmodel.title = '';
        },

        getWjList(){
            this.loading=true;
            API.ShowQuestionaires().then(res=>{
                if(res.code==200){
                    this.$message({
                        type:"success",
                        message:res.msg+'!',
                    })
                    this.WjList  = res.data
                    this.loading=false;
                    this.lookQuestionaire();
                }else{
                    this.$message({
                        type:"error",
                        message:res.msg+'!',
                    })
                }
            }).catch(error =>{
                this.$notify.error({
                    message : error,
                });
            })
        },
        
        UpdateWj(){
            // console.log("ActiveOp",this.ActiveOp)
            this.AddWjmodel = this.NowSelect;
            this.AddQuesShow = true;
        },

        DeleteWj(){
            this.$confirm('确认删除'+this.NowSelect.title+'?(删除后将无法恢复)','删除确认',{
                confirmButtonText:'确认删除',
                cancelButtonText:'取消删除',
                type:'warning'
            }).then(() =>{
                this.loading=true;
                this.DelArg.WjId = this.NowSelect.wjid;
                this.DelArg.QuesId=0;
                console.log(this.DelArg);
                API.DeleteQuestionaire({data:this.DelArg}).then(res=>{
                    if(res.code==200){
                        this.$message({
                            type:"success",
                            message:res.msg+'!',
                        })
                        this.getWjList();
                        this.loading=false;
                        this.ActiveOp=1;
                    }else{
                        this.$message({
                            type:"error",
                            message:res.msg+'!',
                        })
                    }
                }).catch(error=>{
                    this.$notify.error({
                        message:error,
                    })
                })
            }
            )
        },

        OpenQues(index){
            this.ActiveOp = (index+1).toString();
            this.lookQuestionaire();
        },
        
        lookQuestionaire(){
            // console.log(this.NowSelect.wjid)
            // console.log(this.NowSelect.title)
            //console.log(this.NowSelect.desc);
            
            //this.$refs.Design.init(this.NowSelect.WjId,this.NowSelect.title,this.NowSelect.desc);
            this.$refs.design.init(this.NowSelect.wjid,this.NowSelect.title,this.NowSelect.desc);
            // this.$nextTick(()=>{
            //     this.$refs.design.init(this.NowSelect.wjid,this.NowSelect.title,this.NowSelect.desc);
            // })
            
        },
        loginCheck(){
            if(sessionStorage.getItem('username')!=null){
                this.getWjList();
                this.$emit("state");
            }else{
                this.$message({
                    type:'error',
                    message:'请先登录',
                });
                this.$router.push({path:'/login'});
            }
        }
    },
})
</script>
