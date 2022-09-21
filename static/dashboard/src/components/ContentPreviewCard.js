import React, {Component} from "react"
import ImageViewer from "./ImageViewer"
import PhotobookViewer from "./PhotobookViewer"
import VideoPlayer from "./VideoPlayer"



class ContentPreviewCard extends Component{
    constructor(props){
        super(props)

    }

    render(){
        let file = this.props.file
        let title = (this.props.visibility ? this.props.visibility : file.original_filename)
        
        let component = <>
            <img src={this.props.file.url} alt="content preview" />
            </>

        if(this.props.campaign.type === 'embed'){
            component = <VideoPlayer video={file} />
        }
        else if(this.props.campaign.type === 'photobook'){
            let images = [{
                original: file.url,
                thumbnail: file.thumbnail,
            }]
            component = <PhotobookViewer images={images}  />
        } else{
            switch(file.type){
                case 'video':
                case 'audio':
                component = <VideoPlayer video={file} />
                break;
                case 'image':
                    let images = [{
                        original: file.url,
                        thumbnail:  file.thumbnail,
                    }]
                component = <ImageViewer images={images} />
                break;
                case 'other':
                case 'youtube_embed':
                case 'soundcloud_embed':
                component = <h1>Other file</h1>
                break;
            }
        }
        return (
            <div class="single-blog-grid-01 margin-bottom-30">
                <div class="thumb">
                    {component}
                    <div class="news-date">
                        <h5 class="title">{title.substr(0, 25)}</h5>
                    </div>
                </div>
                <div class="content">
                    <h4 class="title"><a href={this.props.link} target="_blank">{this.props.title}</a><i class="fa fa-link"></i> </h4>
                    <p>{this.props.description}</p>
                </div>
            </div>
        )
    }
}

export default ContentPreviewCard;