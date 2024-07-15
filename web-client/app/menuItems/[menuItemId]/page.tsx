import axios, { AxiosError } from 'axios';
import { cookies } from 'next/headers';
import Image from 'next/image';
import { Button } from '@nextui-org/button';
import { ArrowLeftCircleIcon } from '@heroicons/react/24/outline';
import Link from 'next/link';
import { notFound } from 'next/navigation';

interface Dish {
  id: string;
  name: string;
  categoryName: string;
  imageUrl: string;
}

export async function generateMetadata({ params }: { params: any }) {
  try {
    const cookieStore = cookies();
    const token = cookieStore.get('token')!.value;
    const response = await axios.get(
      `${process.env.NEXT_PUBLIC_BACKEND_URL}/api/menuItems/${params.menuItemId}`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );
    const { name } = response.data;
    return {
      title: name,
    };
  } catch (error) {
    if (axios.isAxiosError(error)) {
      const axiosError = error as AxiosError;
      if (axiosError.response?.status === 404) {
        notFound();
      }
    }
    throw error;
  }
}

export default async function MenuItem({ params }: { params: any }) {
  try {
    const cookieStore = cookies();
    const token = cookieStore.get('token')!.value;
    const response = await axios.get(
      `${process.env.NEXT_PUBLIC_BACKEND_URL}/api/menuItems/${params.menuItemId}`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );

    const menuItem = response.data;

    const { name, categoryName, imageUrl } = menuItem;

    return (
      <div className="fixed inset-0 flex flex-col bg-content1 pb-16">
        <div className="relative w-full aspect-square">
          <Image
            src={imageUrl}
            alt={name}
            fill
            className="object-cover"
            priority
          />
          <Link href="/" className="absolute top-4 left-4 z-10">
            <div className="relative p-1 rounded-full bg-black bg-opacity-30 backdrop-blur-sm">
              <ArrowLeftCircleIcon className="h-10 w-10 text-white drop-shadow-lg" />
            </div>
          </Link>
        </div>
        <div className="flex-1 p-4 flex flex-col justify-between">
          <h1 className="text-2xl font-bold text-default-800 mb-2">{name}</h1>
          <p className="text-lg text-default-600 mb-4">{categoryName}</p>
          <div className="mt-auto">
            <Button
              color="primary"
              size="lg"
              className="w-full bg-mainLight dark:bg-mainDark text-white text-lg py-3"
            >
              Add to Cart
            </Button>
          </div>
        </div>
      </div>
    );
  } catch (error) {
    if (axios.isAxiosError(error)) {
      const axiosError = error as AxiosError;
      if (axiosError.response?.status === 404) {
        notFound();
      }
    }
    throw error;
  }
}
