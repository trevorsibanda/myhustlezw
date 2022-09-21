import React, {Component} from "react"

import ImageGallery from 'react-image-gallery';


import 'react-image-gallery/styles/css/image-gallery.css'


class ImageViewer extends Component{
    render(){
        return(
            <ImageGallery showThumbnails={this.props.showThumbnails} slideInterval={15000} showNav={true} autoPlay={true} items={this.props.images} />
        )
    }
}

export default ImageViewer;