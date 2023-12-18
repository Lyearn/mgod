import clsx from 'clsx';
import Heading from '@theme/Heading';

type FeatureItem = {
  icon: string;
  title: string;
  description: string;
};

const FeatureList: FeatureItem[] = [
  {
    icon: 'icon/plug-large.svg',
    title: 'Simplified Integration',
    description: 'Streamline MongoDB interactions with mgod. Reduce redundancy and enhance type safety.'
  },
  {
    icon: 'icon/setup-large.svg',
    title: 'Flexible and Open Source',
    description: 'With an MIT license, mgod gives you flexibility and control over your MongoDB stack. Ensure transparency, avoid unexpected API changes.',
  },
  {
    icon: 'icon/trending-up-large.svg',
    title: 'Continuous Improvement',
    description: 'Dedicated to delivering regular updates, mgod consistently introduces new features and improvements.',
  },
];

function Feature({ icon, title, description }: FeatureItem) {
  return (
    <div className="xs:max-w-[30rem]">
      <img width="40rem" height="40rem" src={icon} alt={title} className="mb-[1.6rem]" aria-hidden="true" />
      <Heading as="h3" className="mb-[0.4rem] heading-xs text-text-paragraph">{title}</Heading>
      <p className="body-short-01 text-text-paragraph">{description}</p>
    </div>
  );
}

export default function HomepageFeatures(): JSX.Element {
  return (
    <section className="sm:py-[10rem] px-[4.8rem] mx-auto flex flex-row flex-wrap gap-[7.2rem] place-content-evenly">
      <Heading as="h2" className="screen-reader-only">Features</Heading>
      {FeatureList.map((props, idx) => (
        <Feature key={idx} {...props} />
      ))}
    </section>
  );
}
