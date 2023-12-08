import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

type FeatureItem = {
  label: string;
  title: string;
  description: string;
};

const FeatureList: FeatureItem[] = [
  {
    label: '01',
    title: 'Simplified Integration',
    description: 'Mgod streamlines MongoDB interactions, reducing redundancy and enhancing type safety, making it a preferred choice for developers.'
  },
  {
    label: '02',
    title: 'Flexible and Open Source',
    description: 'With an MIT license, mgod offers flexibility and control over your MongoDB stack, ensuring transparency and avoiding unexpected API changes.',
  },
  {
    label: '03',
    title: 'Continuous Improvement',
    description: 'We are committed to weekly updates, delivering new features and improvements regularly.',
  },
];

function Feature({label, title, description}: FeatureItem) {
  return (
    <div className={clsx('text--primary', styles.featureItem)}>
      <div>
        <span className="">{label}</span>
      </div>
      <div className="">
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
