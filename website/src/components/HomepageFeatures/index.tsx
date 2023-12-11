import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

type FeatureItem = {
  icon: JSX.Element;
  title: string;
  description: string;
};

const FeatureList: FeatureItem[] = [
  {
    icon: (
      <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path fill-rule="evenodd" clip-rule="evenodd" d="M31.707 9.70215L34.2319 7.17732C34.4271 6.98206 34.4271 6.66548 34.2319 6.47021L33.5248 5.76311C33.3295 5.56785 33.0129 5.56785 32.8177 5.76311L30.2928 8.28799C27.551 6.23252 23.6442 6.45156 21.1507 8.9451L18.3223 11.7735C17.9317 12.1641 17.9317 12.7972 18.3223 13.1877L26.8076 21.673C27.1981 22.0635 27.8312 22.0635 28.2218 21.673L31.0502 18.8446C33.5438 16.3509 33.7628 12.444 31.707 9.70215ZM29.6035 10.3271L29.6682 10.3918C31.5885 12.3469 31.5778 15.4886 29.636 17.4304L27.5147 19.5517L20.4436 12.4806L22.5649 10.3593C24.5067 8.41748 27.6484 8.40675 29.6035 10.3271ZM16.728 17.7071C17.1185 17.3166 17.7517 17.3166 18.1422 17.7071C18.5327 18.0976 18.5327 18.7308 18.1422 19.1213L16.0209 21.2426L18.8493 24.0709L20.9706 21.9495C21.3612 21.559 21.9943 21.559 22.3849 21.9495C22.7754 22.3401 22.7754 22.9732 22.3849 23.3638L20.2635 25.4851L21.6778 26.8994C22.0683 27.2899 22.0683 27.9231 21.6778 28.3136L18.8493 31.1421C16.3558 33.6356 12.4491 33.8546 9.70724 31.7992L7.18236 34.3241C6.9871 34.5193 6.67052 34.5193 6.47526 34.3241L5.76815 33.6169C5.57289 33.4217 5.57289 33.1051 5.76815 32.9098L8.29298 30.385C6.23723 27.6432 6.45618 23.7362 8.94983 21.2426L11.7783 18.4141C12.1688 18.0236 12.8019 18.0236 13.1925 18.4141L14.6067 19.8284L16.728 17.7071ZM13.8905 21.9406L12.4854 20.5355L10.364 22.6568C8.42221 24.5986 8.41148 27.7403 10.3318 29.6954L10.3965 29.76C12.3516 31.6804 15.4933 31.6697 17.4351 29.7278L19.5564 27.6065L13.9087 21.9588L13.8995 21.9497L13.8905 21.9406Z" fill="black" />
      </svg>
    ),
    title: 'Simplified Integration',
    description: 'Mgod streamlines MongoDB interactions, reducing redundancy and enhancing type safety, making it a preferred choice for developers.'
  },
  {
    icon: (
      <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path fill-rule="evenodd" clip-rule="evenodd" d="M29.874 11C29.4299 9.27477 27.8638 8 26 8C24.1362 8 22.5701 9.27477 22.126 11H8C7.44772 11 7 11.4477 7 12C7 12.5523 7.44772 13 8 13H22.126C22.5701 14.7252 24.1362 16 26 16C27.8638 16 29.4299 14.7252 29.874 13H32C32.5523 13 33 12.5523 33 12C33 11.4477 32.5523 11 32 11H29.874ZM28 12C28 13.1046 27.1046 14 26 14C24.8954 14 24 13.1046 24 12C24 10.8954 24.8954 10 26 10C27.1046 10 28 10.8954 28 12Z" fill="black" />
        <path fill-rule="evenodd" clip-rule="evenodd" d="M8 19C7.44772 19 7 19.4477 7 20C7 20.5523 7.44772 21 8 21H9.12602C9.57006 22.7252 11.1362 24 13 24C14.8638 24 16.4299 22.7252 16.874 21H32C32.5523 21 33 20.5523 33 20C33 19.4477 32.5523 19 32 19H16.874C16.4299 17.2748 14.8638 16 13 16C11.1362 16 9.57006 17.2748 9.12602 19H8ZM13 18C14.1046 18 15 18.8954 15 20C15 21.1046 14.1046 22 13 22C11.8954 22 11 21.1046 11 20C11 18.8954 11.8954 18 13 18Z" fill="black" />
        <path fill-rule="evenodd" clip-rule="evenodd" d="M7 28C7 27.4477 7.44772 27 8 27H20.126C20.5701 25.2748 22.1362 24 24 24C25.8638 24 27.4299 25.2748 27.874 27H32C32.5523 27 33 27.4477 33 28C33 28.5523 32.5523 29 32 29H27.874C27.4299 30.7252 25.8638 32 24 32C22.1362 32 20.5701 30.7252 20.126 29H8C7.44772 29 7 28.5523 7 28ZM26 28C26 26.8954 25.1046 26 24 26C22.8954 26 22 26.8954 22 28C22 29.1046 22.8954 30 24 30C25.1046 30 26 29.1046 26 28Z" fill="black" />
      </svg>
    ),
    title: 'Flexible and Open Source',
    description: 'With an MIT license, mgod offers flexibility and control over your MongoDB stack, ensuring transparency and avoiding unexpected API changes.',
  },
  {
    icon: (
      <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path fill-rule="evenodd" clip-rule="evenodd" d="M27.9329 14.64H24.3471C23.7948 14.64 23.3471 14.1923 23.3471 13.64C23.3471 13.0877 23.7948 12.64 24.3471 12.64H30.3471C30.8994 12.64 31.3471 13.0877 31.3471 13.64V19.64C31.3471 20.1923 30.8994 20.64 30.3471 20.64C29.7948 20.64 29.3471 20.1923 29.3471 19.64V16.0542L22.0542 23.3471C21.6976 23.7037 21.1313 23.739 20.7332 23.4294L16.93 20.4713L10.7613 26.64C10.3708 27.0305 9.73765 27.0305 9.34712 26.64C8.9566 26.2495 8.9566 25.6163 9.34712 25.2258L16.14 18.4329C16.4966 18.0763 17.063 18.041 17.4611 18.3506L21.2642 21.3087L27.9329 14.64Z" fill="black" />
      </svg>
    ),
    title: 'Continuous Improvement',
    description: 'We are committed to weekly updates, delivering new features and improvements regularly.',
  },
];

function Feature({ icon, title, description }: FeatureItem) {
  return (
    <div className={clsx('text--primary', styles.featureItem)}>
      <div style={{ marginBlockEnd: 16 }}>
        {icon}
      </div>
      <div>
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): JSX.Element {
  return (
    <section className={clsx(styles.features)}>
      {FeatureList.map((props, idx) => (
        <Feature key={idx} {...props} />
      ))}
    </section>
  );
}
