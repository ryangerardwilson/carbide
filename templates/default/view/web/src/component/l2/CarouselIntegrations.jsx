import { Carousel } from './Carousel.jsx';

export function GlideCarousel(props) {
  return (
    <div data-integration="Glide">
      <Carousel label="Glide carousel" {...props} />
    </div>
  );
}

export function SplideCarousel(props) {
  return (
    <div data-integration="Splide">
      <Carousel label="Splide carousel" {...props} />
    </div>
  );
}
