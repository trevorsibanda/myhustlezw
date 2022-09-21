import React, { Component } from "react"

import ImageGallery from 'react-image-gallery';

import 'react-image-gallery/styles/css/image-gallery.css'

class PhotobookViewer extends Component {
    render() {
        return (
            <ImageGallery showThumbnails={false} autoPlay={true}  items={this.props.images} />
        )
    }
}

export default PhotobookViewer;